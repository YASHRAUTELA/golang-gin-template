package routes

import (
	"server/controllers"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(route *gin.Engine) {
	userRoutes := route.Group("/user")
	userRoutes.Use(middleware.AuthMiddleware())
	{
		userRoutes.GET("/", controllers.GetUser)
	}
}
