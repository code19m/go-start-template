package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (srv *HttpServer) setupApi() {
	r := srv.router

	baseRoute := r.Group("/api/v1/")

	// Register your handlers here
	{
		baseRoute.POST("my-model/", srv.createMyModelHandler)
	}
}

func (srv *HttpServer) setupHealthCheck() {
	srv.router.GET("/health", func(c *gin.Context) {
		c.Status(200)
	})
}

// @title go-start-template API
// @description This document contains the source for the go-start-template API
// @BasePath /api/v1/
func (srv *HttpServer) setupSwaggerDocs() {
	baseRoute := srv.router.Group("/api/v1/")

	ginSwagger.WrapHandler(
		swaggerFiles.Handler,
	)
	baseRoute.GET("/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler))
}
