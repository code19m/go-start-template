package http

import (
	"context"
	"go-start-template/internal/config"
	"go-start-template/internal/domain"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type myModelSrv interface {
	Create(ctx context.Context, params domain.CreateMyModelParams) (int32, error)
	FindOne(ctx context.Context, id int32) (domain.MyModel, error)
}

type HttpServer struct {
	*http.Server

	serverConfig  *config.HttpServer
	openApiConfig *config.OpenAPI
	log           *slog.Logger
	router        *gin.Engine
	myModelSrv    myModelSrv
	addr          string
}

func New(
	srvConfig *config.HttpServer,
	openApiConfig *config.OpenAPI,

	log *slog.Logger,
	appmode string,
	addr string,

	// Services
	myModelSrv myModelSrv,
) (
	*HttpServer, error,
) {
	setEngineMode(appmode)

	router := gin.New()

	srv := &HttpServer{
		serverConfig:  srvConfig,
		openApiConfig: openApiConfig,
		log:           log,
		router:        router,
		myModelSrv:    myModelSrv,
		addr:          addr,

		// Ignore ReadTimeout warning since used http.TimeoutHandler instead
		Server: &http.Server{ //nolint: gosec
			Handler:     http.TimeoutHandler(router, srvConfig.TimeOut, "Server timeout"),
			Addr:        addr,
			IdleTimeout: srvConfig.IdleTimeout,
		},
	}

	srv.router.ContextWithFallback = true
	srv.setupGlobalMiddlewares()
	srv.setupApi()
	srv.registerCustomValidators()

	return srv, nil
}

func setEngineMode(mode string) {
	// Set gin mode to release mod, so we don't need any default logs from gin
	gin.SetMode(gin.ReleaseMode)
}

func (srv *HttpServer) setupGlobalMiddlewares() {
	srv.router.Use(
		accessLoggerMiddleware(srv.log),
		corsMiddleware(),
		gin.Recovery(),
	)
}
