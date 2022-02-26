package internal

import (
	"github.com/theartofdevel/imgur-service/internal/config"
	"github.com/theartofdevel/imgur-service/internal/events"
	imgurService "github.com/theartofdevel/imgur-service/internal/imgur"
	"github.com/theartofdevel/imgur-service/pkg/client/imgur"
	"github.com/theartofdevel/imgur-service/pkg/client/mq/rabbitmq"
	"github.com/theartofdevel/imgur-service/pkg/logging"
	"net/http"
	"sync"
)

type app struct {
	cfg        *config.Config
	logger     *logging.Logger
	httpServer *http.Server
	service    imgurService.Service
}

type App interface {
	Run()
}

func NewApp(logger *logging.Logger, cfg *config.Config) (App, error) {
	client := imgur.NewClient(cfg.Imgur.URL, cfg.Imgur.AccessToken, cfg.Imgur.ClientID, &http.Client{})
	service := imgurService.NewImgurService(client, logger)

	return &app{
		cfg:     cfg,
		logger:  logger,
		service: service,
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

	var wg sync.WaitGroup

	for i := 0; i < a.cfg.AppConfig.EventWorkers; i++ {
		worker := events.NewWorker(i, consumer, a.cfg.RabbitMQ.Producer.Queue, producer, messages, a.logger, a.service, &wg)

		wg.Add(1)
		go worker.Process()
		a.logger.Infof("Event Worker #%d started", i)
	}

	wg.Wait()
}
