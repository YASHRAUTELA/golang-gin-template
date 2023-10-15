package routes

import (
	"myapp/middlewares"

	"github.com/gin-gonic/gin"
)

func InitRoutes(router gin.Engine) {
	// use middleware
	InitMiddleware(router)
	AlbumRoutes(router)
}

func InitMiddleware(router gin.Engine) {
	router.GET("/", func(ctx *gin.Context) {
		ctx.IndentedJSON(200, "Welcome to Gin")
	})
	router.Use(middlewares.CORSMiddleware())
}
