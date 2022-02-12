package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	IsDebug       *bool `yaml:"is_debug" env:"ST-BOT-IsDebug" env-default:"false"  env-required:"true"`
	IsDevelopment *bool `yaml:"is_development" env:"ST-BOT-IsDevelopment" env-default:"false" env-required:"true"`
	Listen        struct {
		Type   string `yaml:"type" env:"ST-BOT-ListenType" env-default:"port"`
		BindIP string `yaml:"bind_ip" env:"ST-BOT-BindIP" env-default:"localhost"`
		Port   string `yaml:"port" env:"ST-BOT-Port" env-default:"8080"`
	} `yaml:"listen" env-required:"true"`
	Telegram struct {
		Token string `yaml:"token" env:"ST-BOT-TelegramToken" env-required:"true"`
	}
	RabbitMQ struct {
		Host     string `yaml:"host" env:"ST-BOT-RabbitHost" env-required:"true"`
		Port     string `yaml:"port" env:"ST-BOT-RabbitPort" env-required:"true"`
		Username string `yaml:"username" env:"ST-BOT-RabbitUsername" env-required:"true"`
		Password string `yaml:"password" env:"ST-BOT-RabbitPassword" env-required:"true"`
		Consumer struct {
			Queue              string `yaml:"queue" env:"ST-BOT-RabbitConsumerQueue" env-required:"true"`
			MessagesBufferSize int    `yaml:"messages_buff_size" env:"ST-BOT-RabbitConsumerMBS" env-default:"100"`
		} `yaml:"consumer" env-required:"true"`
		Producer struct {
			Queue string `yaml:"queue" env:"ST-BOT-RabbitProducerQueue" env-required:"true"`
		} `yaml:"producer" env-required:"true"`
	}
	Imgur struct {
		AccessToken  string `yaml:"access_token" env:"ST-BOT-ImgurAccessToken" env-required:"true"`
		ClientSecret string `yaml:"client_secret" env:"ST-BOT-ImgurClientSecret" env-required:"true"`
		ClientID     string `yaml:"client_id" env:"ST-BOT-ImgurClientID" env-required:"true"`
		URL          string `yaml:"url" env:"ST-BOT-ImgurURL" env-required:"true"`
	} `yaml:"imgur"`
	AppConfig AppConfig `yaml:"app" env-required:"true"`
}

type AppConfig struct {
	EventWorkers int    `yaml:"event_workers" env:"ST-BOT-EventWorks" env-default:"3" env-required:"true"`
	LogLevel     string `yaml:"log_level" env:"ST-BOT-LogLevel" env-default:"error" env-required:"true"`
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
