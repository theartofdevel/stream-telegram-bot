package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/theartofdevel/telegram_bot/internal/config"
	"github.com/theartofdevel/telegram_bot/internal/events"
	"github.com/theartofdevel/telegram_bot/internal/service"
	"github.com/theartofdevel/telegram_bot/pkg/client/imgur"
	"github.com/theartofdevel/telegram_bot/pkg/client/mq"
	"github.com/theartofdevel/telegram_bot/pkg/client/mq/rabbitmq"
	"github.com/theartofdevel/telegram_bot/pkg/logging"
	tele "gopkg.in/telebot.v3"
	"net/http"
	"time"
)

type app struct {
	cfg          *config.Config
	logger       *logging.Logger
	httpServer   *http.Server
	imgurService service.ImgurService
	bot          *tele.Bot
	producer     mq.Producer
}

type App interface {
	Run()
}

func NewApp(logger *logging.Logger, cfg *config.Config) (App, error) {

	client := http.Client{}
	imgurClient := imgur.NewClient(cfg.Imgur.URL, cfg.Imgur.AccessToken, cfg.Imgur.ClientID, &client)
	imgurService := service.NewImgurService(imgurClient, logger)

	return &app{
		cfg:          cfg,
		logger:       logger,
		imgurService: imgurService,
	}, nil
}

func (a *app) Run() {
	a.startBot()
	a.startConsume()
	// TODO fixMe
	a.bot.Start()
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

	for i := 0; i < a.cfg.AppConfig.EventWorkers; i++ {
		worker := events.NewWorker(i, consumer, a.bot, producer, messages, a.logger)

		go worker.Process()
		a.logger.Infof("Event Worker #%d started", i)
	}

	a.producer = producer
}

func (a *app) startBot() {
	pref := tele.Settings{
		Token:   a.cfg.Telegram.Token,
		Poller:  &tele.LongPoller{Timeout: 60 * time.Second},
		Verbose: false,
		OnError: a.OnBotError,
	}
	var botErr error
	a.bot, botErr = tele.NewBot(pref)
	if botErr != nil {
		a.logger.Fatal(botErr)
		return
	}

	a.bot.Handle("/yt", func(c tele.Context) error {
		trackName := c.Message().Payload

		request := events.SearchTrackRequest{
			RequestID: fmt.Sprintf("%d", c.Sender().ID),
			Name:      trackName,
		}

		marshal, _ := json.Marshal(request)

		err := a.producer.Publish(a.cfg.RabbitMQ.Producer.Queue, marshal)
		if err != nil {
			return c.Send(fmt.Sprintf("ошибка: %s", err.Error()))
		}

		return c.Send(fmt.Sprintf("Заявка принята"))
	})

	a.bot.Handle(tele.OnPhoto, func(c tele.Context) error {
		// Photos only.
		photo := c.Message().Photo
		file, err := a.bot.File(&photo.File)
		if err != nil {
			return c.Send("Не удалось скачать изображение")
		}
		defer file.Close()
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(file)
		if err != nil {
			return c.Send("Не удалось скачать изображение")
		}

		if buf.Len() > 10_485_760 {
			return c.Send("Лимит 10МБ")
		}

		image, err := a.imgurService.ShareImage(context.Background(), buf.Bytes())
		if err != nil {
			return c.Send("Не удалось залить изображение")
		}

		return c.Send(image)
	})
}

func (a *app) OnBotError(err error, ctx tele.Context) {
	a.logger.Error(err)
}
