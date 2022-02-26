package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	IsDebug       *bool `yaml:"is_debug" env:"IS_IS_DEBUG" env-default:"false"  env-required:"true"`
	IsDevelopment *bool `yaml:"is_development" env:"IS_IS_DEVELOPMENT" env-default:"false" env-required:"true"`
	RabbitMQ      struct {
		Host     string `yaml:"host" env:"IS_RABBIT_HOST" env-required:"true"`
		Port     string `yaml:"port" env:"IS_RABBIT_PORT" env-required:"true"`
		Username string `yaml:"username" env:"IS_RABBIT_USERNAME" env-required:"true"`
		Password string `yaml:"password" env:"IS_RABBIT_PASSWORD" env-required:"true"`
		Consumer struct {
			Queue              string `yaml:"queue" env:"IS_RABBIT_CONSUMER_QUEUE" env-required:"true"`
			MessagesBufferSize int    `yaml:"messages_buff_size" env:"IS_RABBIT_CONSUMER_MBS" env-default:"100"`
		} `yaml:"consumer" env-required:"true"`
		Producer struct {
			Queue string `yaml:"queue" env:"IS_RABBIT_PRODUCER_QUEUE" env-required:"true"`
		} `yaml:"producer" env-required:"true"`
	}
	Imgur struct {
		AccessToken  string `yaml:"access_token" env:"IS_IMGUR_ACCESS_TOKEN" env-required:"true"`
		ClientSecret string `yaml:"client_secret" env:"IS_IMGUR_CLIENT_SECRET" env-required:"true"`
		ClientID     string `yaml:"client_id" env:"IS_IMGUR_CLIENT_ID" env-required:"true"`
		URL          string `yaml:"url" env:"IS_IMGUR_URL" env-required:"true"`
	} `yaml:"imgur"`
	AppConfig AppConfig `yaml:"app" env-required:"true"`
}

type AppConfig struct {
	EventWorkers int    `yaml:"event_workers" env:"IS_EVENTWORKS" env-default:"3" env-required:"true"`
	LogLevel     string `yaml:"log_level" env:"IS_LOGLEVEL" env-default:"error" env-required:"true"`
}

var instance *Config
var once sync.Once

func GetConfig(path string) *Config {
	once.Do(func() {
		log.Printf("read application config in path %s", path)

		instance = &Config{}

		if err := cleanenv.ReadConfig(path, instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			log.Print(help)
			log.Fatal(err)
		}
	})
	return instance
}
