package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	IsDebug       bool `env:"YTS_IS_DEBUG" env-default:"false"`
	IsDevelopment bool `env:"YTS_IS_DEVELOPMENT" env-default:"false"`
	YouTube       struct {
		APIURL          string `env:"YTS_YT_API_URL" env-required:"true"`
		RefreshTokenURL string `env:"YTS_YT_RefreshTokenURL" env-required:"true"`
		APIKey          string `env:"YTS_YT_APIKey" env-required:"true"`
		ClientID        string `env:"YTS_YT_CLIENT_ID" env-required:"true"`
		ClientSecret    string `env:"YTS_YT_CLIENT_SECRET" env-required:"true"`
		AccessToken     string `env:"YTS_YT_ACCESS_TOKEN" env-required:"true"`
		RefreshToken    string `env:"YTS_YT_REFRESH_TOKEN" env-required:"true"`
		AuthRedirectUri string `env:"YTS_YT_AUTH_REDIRECT_URI" env-required:"true"`
		AuthSuccessUri  string `env:"YTS_YT_AUTH_SUCCESS_URI" env-required:"true"`
		AccountsUri     string `env:"YTS_YT_ACCOUNTS_URI" env-required:"true"`
	}
	RabbitMQ struct {
		Host     string `env:"YTS_RABBIT_HOST" env-required:"true"`
		Port     string `env:"YTS_RABBIT_PORT" env-required:"true"`
		Username string `env:"YTS_RABBIT_USERNAME" env-required:"true"`
		Password string `env:"YTS_RABBIT_PASSWORD" env-required:"true"`
		Consumer struct {
			Queue              string `env:"YTS_RABBIT_CONSUMER_QUEUE" env-required:"true"`
			MessagesBufferSize int    `env:"YTS_RABBIT_CONSUMER_MBS" env-default:"100"`
		}
		Producer struct {
			Queue string `env:"YTS_Rabbit_PRODUCERQUEUE" env-required:"true"`
		}
	}
	AppConfig AppConfig
}

type AppConfig struct {
	EventWorkers int    `env:"YTS_EVENT_WORKERS" env-default:"3"`
	LogLevel     string `env:"YTS_LOG_LEVEL" env-default:"error"`
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
