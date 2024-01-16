package http

import (
	"context"
	"go-start-template/internal/config"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpHandler struct {
	serverConfig  *config.HttpServer
	openApiConfig *config.OpenAPI
	log           *slog.Logger
	server        *http.Server
	router        *gin.Engine
	addr          string
}

func New(
	srvConfig *config.HttpServer,
	openApiConfig *config.OpenAPI,

	log *slog.Logger,
	// services *service.Pack,
	// clients *client.Pack,
	appmode string,
	addr string,
) (
	*HttpHandler, error,
) {
	setEngineMode(appmode)

	router := gin.New()

	srv := &HttpHandler{
		serverConfig:  srvConfig,
		openApiConfig: openApiConfig,
		log:           log,
		router:        router,
		// services: services,
		// clients:  clients,
		addr: addr,

		// Ignore ReadTimeout warning since used http.TimeoutHandler instead
		server: &http.Server{ //nolint: gosec
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

func (srv *HttpHandler) Run() error {
	return srv.server.ListenAndServe()
}

func (srv *HttpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	srv.router.ServeHTTP(w, req)
}

func (srv *HttpHandler) Shutdown(ctx context.Context) error {
	return srv.server.Shutdown(ctx)
}

func setEngineMode(mode string) {
	ginMode := gin.ReleaseMode
	switch mode {
	case config.LocalMode:
		ginMode = gin.DebugMode
	case config.TestMode:
		ginMode = gin.TestMode
	}
	gin.SetMode(ginMode)
}

func (srv *HttpHandler) setupGlobalMiddlewares() {
	srv.router.Use(
		accessLoggerMiddleware(srv.log),
		corsMiddleware(),
		gin.Recovery(),
	)
}
