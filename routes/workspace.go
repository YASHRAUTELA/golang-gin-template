package routes

import (
	"server/controllers"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

func WorkspaceRoutes(route *gin.Engine) {
	workspaceRoutes := route.Group("/workspace")
	workspaceRoutes.Use(middleware.AuthMiddleware())
	{
		workspaceRoutes.POST("/", controllers.CreateWorkspace)
	}
}
