package routes

import (
	"server/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(route *gin.Engine) {
	authRoutes := route.Group("/auth")
	{
		authRoutes.POST("/login", controllers.Login)
		authRoutes.POST("/register", controllers.Register)
		authRoutes.POST("/sociallogin", controllers.SocialLogin)
	}
}
