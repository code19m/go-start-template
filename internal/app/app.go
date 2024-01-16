package app

import (
	"go-start-template/internal/config"
	"go-start-template/internal/handler/http"
	"go-start-template/pkg/logger"
	"log"
)

func Run(http_addr string) {
	// Initilize config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %s", err.Error())
	}

	// Initilize logger
	logger, err := logger.NewSlogLogger(cfg.Logger.Level, cfg.Logger.Format)
	if err != nil {
		log.Fatalf("failed to load logger: %s", err.Error())
	}
	logger.Debug("Initialized logger")

	// TODO: Init repository

	// TODO: Init service

	httpHandler, err := http.New(&cfg.HttpServer, &cfg.OpenAPI, logger, cfg.AppMode, http_addr)
	if err != nil {
		log.Fatalf("failed to initialize httpHandler: %s", err.Error())
	}

	httpHandler.Run()

	// TODO: Init kafka consumer

	// TODO: Graceful shutdown

}

// TODO: Implement Hot-reload, Postgres connection, Goose migrations, Makefile, Dockerfile, Docker-compose
