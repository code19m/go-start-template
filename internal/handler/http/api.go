package http

import (
	"fmt"
	"go-start-template/api/gen/openapi"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (srv *HttpHandler) setupApi() {
	r := srv.router
	r.GET("/health", checkHealth)

	baseRoute := r.Group("/api/v1/")

	openApiAddr := fmt.Sprintf("%s:%d", srv.openApiConfig.Host, srv.openApiConfig.Port)
	SetupSwaggerDocs(baseRoute, openApiAddr)

	// Register your handlers here
	{
		baseRoute.GET("/hello", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "Hello World!",
			})
		})
	}
}

func checkHealth(c *gin.Context) {
	c.Status(200)
}

// @title go-start-template API
// @description This document contains the source for the go-start-template API
// @BasePath /api/v1/
func SetupSwaggerDocs(baseRoute *gin.RouterGroup, addr string) {
	ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL(fmt.Sprintf(
			"%s/swagger/docs.json",
			addr,
		)),
	)
	openapi.SwaggerInfo.Host = addr
	baseRoute.GET("/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler))
}
