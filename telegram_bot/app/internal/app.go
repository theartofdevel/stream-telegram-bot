package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/theartofdevel/telegram_bot/internal/config"
	"github.com/theartofdevel/telegram_bot/internal/events"
	"github.com/theartofdevel/telegram_bot/internal/events/imgur"
	"github.com/theartofdevel/telegram_bot/internal/events/youtube"
	"github.com/theartofdevel/telegram_bot/internal/service/bot"
	"github.com/theartofdevel/telegram_bot/pkg/client/mq"
	"github.com/theartofdevel/telegram_bot/pkg/client/mq/rabbitmq"
	"github.com/theartofdevel/telegram_bot/pkg/logging"
	tele "gopkg.in/telebot.v3"
	"net/http"
	"time"
)

type app struct {
	cfg                    *config.Config
	logger                 *logging.Logger
	httpServer             *http.Server
	producer               mq.Producer
	youtubeProcessStrategy events.ProcessEventStrategy
	imgurProcessStrategy   events.ProcessEventStrategy
	bot                    *tele.Bot
}

type App interface {
	Run()
}

func NewApp(logger *logging.Logger, cfg *config.Config) (App, error) {
	return &app{
		cfg:                    cfg,
		logger:                 logger,
		youtubeProcessStrategy: youtube.NewYouTubeProcessEventStrategy(logger),
		imgurProcessStrategy:   imgur.NewImgurProcessEventStrategy(logger),
	}, nil
}

func (a *app) Run() {
	bot, err := a.createBot()
	if err != nil {
		return
	}
	a.bot = bot
	a.startConsume()
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

	err = consumer.DeclareQueue(a.cfg.RabbitMQ.Consumer.Youtube, true, false, false, nil)
	if err != nil {
		a.logger.Fatal(err)
	}
	ytMessages, err := consumer.Consume(a.cfg.RabbitMQ.Consumer.Youtube)
	if err != nil {
		a.logger.Fatal(err)
	}

	botService := bot.Service{
		Bot:    a.bot,
		Logger: a.logger,
	}

	for i := 0; i < a.cfg.AppConfig.EventWorkers.Youtube; i++ {
		worker := events.NewWorker(i, consumer, a.youtubeProcessStrategy, botService, producer, ytMessages, a.logger)

		go worker.Process()
		a.logger.Infof("YouTube Event Worker #%d started", i)
	}

	err = consumer.DeclareQueue(a.cfg.RabbitMQ.Consumer.Imgur, true, false, false, nil)
	if err != nil {
		a.logger.Fatal(err)
	}
	imgurMessages, err := consumer.Consume(a.cfg.RabbitMQ.Consumer.Imgur)
	if err != nil {
		a.logger.Fatal(err)
	}

	for i := 0; i < a.cfg.AppConfig.EventWorkers.Imgur; i++ {
		worker := events.NewWorker(i, consumer, a.imgurProcessStrategy, botService, producer, imgurMessages, a.logger)

		go worker.Process()
		a.logger.Infof("Imgur Event Worker #%d started", i)
	}

	a.producer = producer
}

func (a *app) createBot() (abot *tele.Bot, botErr error) {
	pref := tele.Settings{
		Token:   a.cfg.Telegram.Token,
		Poller:  &tele.LongPoller{Timeout: 60 * time.Second},
		Verbose: false,
		OnError: a.OnBotError,
	}
	abot, botErr = tele.NewBot(pref)
	if botErr != nil {
		a.logger.Fatal(botErr)
		return
	}

	abot.Handle("/start", func(c tele.Context) error {
		return c.Send(fmt.Sprintf("/yt - find youtube track by name\nupload photo with compressions and get imgur short url"))
	})

	abot.Handle("/help", func(c tele.Context) error {
		return c.Send(fmt.Sprintf("/yt - find youtube track by name\nupload photo with compressions and get imgur short url"))
	})

	abot.Handle("/yt", func(c tele.Context) error {
		trackName := c.Message().Payload

		request := youtube.SearchTrackRequest{
			RequestID: fmt.Sprintf("%d", c.Sender().ID),
			Name:      trackName,
		}

		marshal, _ := json.Marshal(request)

		err := a.producer.Publish(a.cfg.RabbitMQ.Producer.Youtube, marshal)
		if err != nil {
			return c.Send(fmt.Sprintf("ошибка: %s", err.Error()))
		}

		return c.Send(fmt.Sprintf("Заявка принята"))
	})

	abot.Handle(tele.OnPhoto, func(c tele.Context) error {
		// Photos only.
		photo := c.Message().Photo
		file, err := abot.File(&photo.File)
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

		request := imgur.UploadImageRequest{
			RequestID: fmt.Sprintf("%d", c.Sender().ID),
			Photo:     buf.Bytes(),
		}

		marshal, _ := json.Marshal(request)

		err = a.producer.Publish(a.cfg.RabbitMQ.Producer.Imgur, marshal)
		if err != nil {
			return c.Send(fmt.Sprintf("ошибка: %s", err.Error()))
		}

		return c.Send(fmt.Sprintf("Заявка принята"))
	})

	return
}

func (a *app) OnBotError(err error, ctx tele.Context) {
	a.logger.Error(err)
}
