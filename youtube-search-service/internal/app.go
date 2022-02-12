package internal

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/theartofdevel/youtube-search-service/internal/config"
	"github.com/theartofdevel/youtube-search-service/internal/events"
	youtubeService "github.com/theartofdevel/youtube-search-service/internal/youtube"
	"github.com/theartofdevel/youtube-search-service/pkg/client/mq/rabbitmq"
	"github.com/theartofdevel/youtube-search-service/pkg/client/youtube"
	"github.com/theartofdevel/youtube-search-service/pkg/logging"
	"github.com/theartofdevel/youtube-search-service/pkg/shutdown"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"syscall"
	"time"
)

type app struct {
	cfg            *config.Config
	logger         *logging.Logger
	httpServer     *http.Server
	youtubeService youtubeService.Service
	router         *httprouter.Router
}

type App interface {
	Run()
}

func NewApp(logger *logging.Logger, cfg *config.Config) (App, error) {
	logger.Println("router initializing")
	router := httprouter.New()

	ytClient := youtube.NewClient(cfg.Youtube.APIURL, cfg.Youtube.AccessToken, &http.Client{})
	yts := youtubeService.NewService(ytClient, logger)

	return &app{
		cfg:            cfg,
		logger:         logger,
		youtubeService: yts,
		router:         router,
	}, nil
}

func (a *app) Run() {
	a.startConsume()
	a.startHTTP()
}

func (a *app) startConsume() {
	a.logger.Info("start consuming")

	consumer, err := rabbitmq.NewRabbitMQConsumer(rabbitmq.ConsumerConfig{
		BaseConfig: rabbitmq.BaseConfig{
			Host:     a.cfg.RabbitMQ.Host,
			Port:     a.cfg.RabbitMQ.Port,
			Username: a.cfg.RabbitMQ.Username,
			Password: a.cfg.RabbitMQ.Password,
		},
		PrefetchCount: a.cfg.RabbitMQ.Consumer.MessagesBufferSize,
	})
	if err != nil {
		a.logger.Fatal(err)
	}
	producer, err := rabbitmq.NewRabbitMQProducer(rabbitmq.ProducerConfig{
		BaseConfig: rabbitmq.BaseConfig{
			Host:     a.cfg.RabbitMQ.Host,
			Port:     a.cfg.RabbitMQ.Port,
			Username: a.cfg.RabbitMQ.Username,
			Password: a.cfg.RabbitMQ.Password,
		},
	})
	if err != nil {
		a.logger.Fatal(err)
	}

	messages, err := consumer.Consume(a.cfg.RabbitMQ.Consumer.Queue)
	if err != nil {
		a.logger.Fatal(err)
	}

	client := http.Client{}
	ytClient := youtube.NewClient(a.cfg.Youtube.APIURL, a.cfg.Youtube.AccessToken, &client)
	ytService := youtubeService.NewService(ytClient, a.logger)

	for i := 0; i < a.cfg.AppConfig.EventWorkers; i++ {
		worker := events.NewWorker(i, consumer, a.cfg.RabbitMQ.Producer.Queue, producer, messages, a.logger, ytService)

		go worker.Process()
		a.logger.Infof("Event Worker #%d started", i)
	}
}

func (a *app) startHTTP() {
	a.logger.Info("start HTTP")

	var listener net.Listener

	if a.cfg.Listen.Type == "sock" {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			a.logger.Fatal(err)
		}
		socketPath := path.Join(appDir, "app.sock")
		a.logger.Infof("socket path: %s", socketPath)

		a.logger.Info("create and listen unix socket")
		listener, err = net.Listen("unix", socketPath)
		if err != nil {
			a.logger.Fatal(err)
		}
	} else {
		a.logger.Infof("bind application to host: %s and port: %s", a.cfg.Listen.BindIP, a.cfg.Listen.Port)
		var err error
		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", a.cfg.Listen.BindIP, a.cfg.Listen.Port))
		if err != nil {
			a.logger.Fatal(err)
		}
	}

	c := cors.New(cors.Options{
		AllowedMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodOptions, http.MethodDelete},
		AllowedOrigins:     []string{"*"},
		AllowCredentials:   true,
		AllowedHeaders:     []string{"Access-Token", "Refresh-Token", "Authorization", "Location", "Charset", "Access-Control-Allow-Origin", "Content-Type", "content-type", "Origin", "Accept", "Content-Length", "Accept-Encoding", "X-CSRF-Token"},
		OptionsPassthrough: true,
		ExposedHeaders:     []string{"Access-Token", "Refresh-Token", "Location", "Authorization", "Content-Disposition"},
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})

	handler := c.Handler(a.router)

	a.httpServer = &http.Server{
		Handler:      handler,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go shutdown.Graceful([]os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM},
		a.httpServer)

	a.logger.Println("application completely initialized and started")

	if err := a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			a.logger.Warn("server shutdown")
		default:
			a.logger.Fatal(err)
		}
	}
}
