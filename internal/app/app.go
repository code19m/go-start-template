package app

import (
	"context"
	"errors"
	"go-start-template/internal/config"
	httpServer "go-start-template/internal/handler/http"
	"go-start-template/internal/repository/postgres"
	"go-start-template/internal/service"
	"go-start-template/pkg/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func Run(http_addr string) {
	// Load config
	start := time.Now()
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %s\n", err.Error()) // Using standard logger because logger is not initialized yet
	}
	log.Printf("Loaded config in %s\n", time.Since(start).String()) // Using standard logger because logger is not initialized yet

	// Initialize logger
	start = time.Now()
	logger, err := logger.NewSlogLogger(cfg.Logger.Level, cfg.Logger.Format)
	if err != nil {
		log.Fatalf("Failed to load logger: %s\n", err.Error()) // Using standard logger because logger is not initialized yet
	}
	logger.Info("Initialized logger", "elapsed_time", time.Since(start).String())

	// Initialize pgxpool.Pool
	start = time.Now()
	pool, err := postgres.NewConnPool(&cfg.Postgres)
	if err != nil {
		logger.Error("Failed to initialize connection pool", "error", err.Error())
		os.Exit(1)
	}
	logger.Info("Initialized connection pool", "elapsed_time", time.Since(start).String())

	// Initialize repositories
	start = time.Now()
	myModelStore := postgres.NewMyModelStore(logger, pool)
	// More repositories...
	logger.Info("Initialized repositories", "elapsed_time", time.Since(start).String())

	// Initialize services
	start = time.Now()
	myModelSrv := service.NewMyModelSrv(logger, myModelStore)
	// More services...
	logger.Info("Initialized services", "elapsed_time", time.Since(start).String())

	// Initialize http Server
	start = time.Now()
	httpSrv, err := httpServer.New(&cfg.HttpServer, &cfg.OpenAPI, logger, cfg.AppMode, http_addr, myModelSrv)
	if err != nil {
		logger.Error("Failed to initialize httpServer", "error", err.Error())
		os.Exit(1)
	}
	logger.Info("Initialized httpServer", "elapsed_time", time.Since(start).String())

	go func() {
		err := httpSrv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Error occurred while running http server", "error", err.Error())
			os.Exit(1)
		}
	}()

	// TODO: Init kafka consumer

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	// Set maximum shutdown time to Http server's MaxShutdownTime
	var timeout = cfg.HttpServer.MaxShutdownTime
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var wg sync.WaitGroup

	// Stop receiving http requests from clients and give time to process current requests
	wg.Add(1)
	go func() {
		defer wg.Done()
		start := time.Now()
		err := httpSrv.Shutdown(ctx)
		if err != nil {
			logger.Error("Failed to gracefully shutdown http Server", "error", err.Error())
		} else {
			logger.Info("Gracefully shutdown http Server", "elapsed_time", time.Since(start).String())
		}
	}()

	// Wait for shutdown of all upstream services
	// Then close all downstream services
	wg.Wait()
	pool.Close()

	logger.Info("Application shut down...")
}
