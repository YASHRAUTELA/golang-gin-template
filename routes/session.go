package routes

import (
	"server/controllers"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

func SessionRoutes(route *gin.Engine) {
	sessionRoutes := route.Group("/session")
	sessionRoutes.Use(middleware.AuthMiddleware())
	{
		sessionRoutes.GET("/", controllers.GetUserSessions)
		sessionRoutes.GET("/:sessionId", controllers.GetSession)
		sessionRoutes.POST("/", controllers.CreateSession)
		sessionRoutes.GET("/:sessionId/files", controllers.GetSessionFiles)
		sessionRoutes.POST("/files", controllers.AddSessionFiles)
	}
}
