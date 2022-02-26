package internal

import (
	"github.com/theartofdevel/youtube-search-service/internal/config"
	"github.com/theartofdevel/youtube-search-service/internal/events"
	youtubeService "github.com/theartofdevel/youtube-search-service/internal/youtube"
	"github.com/theartofdevel/youtube-search-service/pkg/client/mq/rabbitmq"
	"github.com/theartofdevel/youtube-search-service/pkg/client/youtube"
	"github.com/theartofdevel/youtube-search-service/pkg/logging"
	"net/http"
	"sync"
)

type app struct {
	cfg            *config.Config
	logger         *logging.Logger
	httpServer     *http.Server
	youtubeService *youtubeService.Service
}

type App interface {
	Run()
}

func NewApp(logger *logging.Logger, cfg *config.Config) (App, error) {
	ytClient := youtube.New(cfg.YouTube.APIURL, cfg.YouTube.AccessToken, cfg.YouTube.RefreshToken, cfg.YouTube.ClientID,
		cfg.YouTube.ClientSecret, cfg.YouTube.RefreshTokenURL, cfg.YouTube.AuthRedirectUri, cfg.YouTube.AccountsUri, nil)
	yts := youtubeService.NewService(ytClient, logger)

	return &app{
		cfg:            cfg,
		logger:         logger,
		youtubeService: yts,
	}, nil
}

func (a *app) Run() {
	a.startConsume()
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

	err = consumer.DeclareQueue(a.cfg.RabbitMQ.Consumer.Queue, true, false, false, nil)
	if err != nil {
		a.logger.Fatal(err)
	}
	messages, err := consumer.Consume(a.cfg.RabbitMQ.Consumer.Queue)
	if err != nil {
		a.logger.Fatal(err)
	}

	wg := sync.WaitGroup{}

	for i := 0; i < a.cfg.AppConfig.EventWorkers; i++ {
		worker := events.NewWorker(i, consumer, a.cfg.RabbitMQ.Producer.Queue, producer, messages, a.logger,
			a.youtubeService, &wg)

		wg.Add(1)
		go worker.Process()
		a.logger.Infof("Event Worker #%d started", i)
	}

	wg.Wait()
}
