package routes

import (
	"net/http"
	"server/config"

	"github.com/gin-gonic/gin"
)

func DefaultRoutes(route *gin.Engine) {
	route.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Route Not Found"})
	})

	route.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	route.GET("/swagger/*any", config.SwaggerConfig())
}
