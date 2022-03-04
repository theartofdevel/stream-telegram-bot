package main

import (
	"github.com/theartofdevel/imgur-service/internal"
	"github.com/theartofdevel/imgur-service/internal/config"
	"github.com/theartofdevel/imgur-service/pkg/logging"
	"log"
)

func main() {
	log.Print("config initializing")
	cfg := config.GetConfig()

	log.Print("logger initializing")
	logging.Init(cfg.AppConfig.LogLevel)
	logger := logging.GetLogger()

	logger.Println("Creating Application")
	app, err := internal.NewApp(logger, cfg)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Println("Running Application")
	app.Run()
}
