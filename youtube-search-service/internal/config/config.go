package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	IsDebug       *bool `yaml:"is_debug" env:"YTS-IsDebug" env-default:"false"  env-required:"true"`
	IsDevelopment *bool `yaml:"is_development" env:"YTS-IsDevelopment" env-default:"false" env-required:"true"`
	Listen        struct {
		Type   string `yaml:"type" env:"YTS-ListenType" env-default:"port"`
		BindIP string `yaml:"bind_ip" env:"YTS-BindIP" env-default:"localhost"`
		Port   string `yaml:"port" env:"YTS-Port" env-default:"8080"`
	} `yaml:"listen" env-required:"true"`
	Youtube struct {
		APIURL      string `yaml:"api_url"`
		AccessToken string `yaml:"access_token"`
	}
	RabbitMQ struct {
		Host     string `yaml:"host" env:"YTS-RabbitHost" env-required:"true"`
		Port     string `yaml:"port" env:"YTS-RabbitPort" env-required:"true"`
		Username string `yaml:"username" env:"YTS-RabbitUsername" env-required:"true"`
		Password string `yaml:"password" env:"YTS-RabbitPassword" env-required:"true"`
		Consumer struct {
			Queue              string `yaml:"queue" env:"YTS-RabbitConsumerQueue" env-required:"true"`
			MessagesBufferSize int    `yaml:"messages_buff_size" env:"YTS-RabbitConsumerMBS" env-default:"100"`
		} `yaml:"consumer" env-required:"true"`
		Producer struct {
			Queue string `yaml:"queue" env:"YTS-RabbitProducerQueue" env-required:"true"`
		} `yaml:"producer" env-required:"true"`
	}
	AppConfig AppConfig `yaml:"app" env-required:"true"`
}

type AppConfig struct {
	EventWorkers int    `yaml:"event_workers" env:"YTS-EventWorks" env-default:"3" env-required:"true"`
	LogLevel     string `yaml:"log_level" env:"YTS-LogLevel" env-default:"error" env-required:"true"`
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
