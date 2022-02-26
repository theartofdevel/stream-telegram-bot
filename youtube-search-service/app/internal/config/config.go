package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	IsDebug       *bool `yaml:"is_debug" env:"YTS_IS_DEBUG" env-default:"false"  env-required:"true"`
	IsDevelopment *bool `yaml:"is_development" env:"YTS_IS_DEVELOPMENT" env-default:"false" env-required:"true"`
	YouTube       struct {
		APIURL          string `yaml:"api_url" env:"YTS_YT_API_URL" env-required:"true"`
		RefreshTokenURL string `yaml:"refresh_token_url" env:"YTS_YT_RefreshTokenURL" env-required:"true"`
		APIKey          string `yaml:"api_key" env:"YTS_YT_APIKey" env-required:"true"`
		ClientID        string `yaml:"client_id" env:"YTS_YT_CLIENT_ID" env-required:"true"`
		ClientSecret    string `yaml:"client_secret" env:"YTS_YT_CLIENT_SECRET" env-required:"true"`
		AccessToken     string `yaml:"access_token" env:"YTS_YT_ACCESS_TOKEN" env-required:"true"`
		RefreshToken    string `yaml:"refresh_token" env:"YTS_YT_REFRESH_TOKEN" env-required:"true"`
		AuthRedirectUri string `yaml:"auth_redirect_uri" env:"YTS_YT_AUTH_REDIRECT_URI" env-required:"true"`
		AuthSuccessUri  string `yaml:"auth_success_uri" env:"YTS_YT_AUTH_SUCCESS_URI" env-required:"true"`
		AccountsUri     string `yaml:"accounts_uri" env:"YTS_YT_ACCOUNTS_URI" env-required:"true"`
	} `yaml:"youtube"`
	RabbitMQ struct {
		Host     string `yaml:"host" env:"YTS_RABBIT_HOST" env-required:"true"`
		Port     string `yaml:"port" env:"YTS_RABBIT_PORT" env-required:"true"`
		Username string `yaml:"username" env:"YTS_RABBIT_USERNAME" env-required:"true"`
		Password string `yaml:"password" env:"YTS_RABBIT_PASSWORD" env-required:"true"`
		Consumer struct {
			Queue              string `yaml:"queue" env:"YTS_RABBIT_CONSUMER_QUEUE" env-required:"true"`
			MessagesBufferSize int    `yaml:"messages_buff_size" env:"YTS_RABBIT_CONSUMER_MBS" env-default:"100"`
		} `yaml:"consumer"`
		Producer struct {
			Queue string `yaml:"queue" env:"YTS_Rabbit_PRODUCERQUEUE" env-required:"true"`
		} `yaml:"producer"`
	}
	AppConfig AppConfig `yaml:"app"`
}

type AppConfig struct {
	EventWorkers int    `yaml:"event_workers" env:"YTS_EVENT_WORKERS" env-default:"3" env-required:"true"`
	LogLevel     string `yaml:"log_level" env:"YTS_LOG_LEVEL" env-default:"error" env-required:"true"`
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
