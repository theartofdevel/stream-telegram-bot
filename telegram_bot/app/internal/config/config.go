package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	IsDebug       *bool `yaml:"is_debug" env:"ST_BOT_IS_DEBUG" env-default:"false"  env-required:"true"`
	IsDevelopment *bool `yaml:"is_development" env:"ST_BOT_IS_DEVELOPMENT" env-default:"false" env-required:"true"`
	Telegram      struct {
		Token string `yaml:"token" env:"ST_BOT_TELEGRAM_TOKEN" env-required:"true"`
	}
	RabbitMQ struct {
		Host     string `yaml:"host" env:"ST_BOT_RABBIT_HOST" env-required:"true"`
		Port     string `yaml:"port" env:"ST_BOT_RABBIT_PORT" env-required:"true"`
		Username string `yaml:"username" env:"ST_BOT_RABBIT_USERNAME" env-required:"true"`
		Password string `yaml:"password" env:"ST_BOT_RABBIT_PASSWORD" env-required:"true"`
		Consumer struct {
			Youtube            string `yaml:"youtube" env:"ST_BOT_RABBIT_CONSUMER_YOUTUBE" env-required:"true"`
			Imgur              string `yaml:"imgur" env:"ST_BOT_RABBIT_CONSUMER_IMGUR" env-required:"true"`
			MessagesBufferSize int    `yaml:"messages_buff_size" env:"ST_BOT_RABBIT_CONSUMER_MBS" env-default:"100"`
		} `yaml:"consumer" env-required:"true"`
		Producer struct {
			Youtube string `yaml:"youtube" env:"ST_BOT_RABBIT_PRODUCER_YOUTUBE" env-required:"true"`
			Imgur   string `yaml:"imgur" env:"ST_BOT_RABBIT_PRODUCER_IMGUR" env-required:"true"`
		} `yaml:"producer" env-required:"true"`
	}
	AppConfig AppConfig `yaml:"app" env-required:"true"`
}

type AppConfig struct {
	EventWorkers struct {
		Youtube int `yaml:"youtube" env:"ST_BOT_EVENT_WORKERS_YT" env-default:"3" env-required:"true"`
		Imgur   int `yaml:"imgur" env:"ST_BOT_EVENT_WORKERS_IMGUR" env-default:"3" env-required:"true"`
	} `yaml:"event_workers"`
	LogLevel string `yaml:"log_level" env:"ST_BOT_LOG_LEVEL" env-default:"error" env-required:"true"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}

		if err := cleanenv.ReadEnv(instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			log.Print(help)
			log.Fatal(err)
		}
	})
	return instance
}
