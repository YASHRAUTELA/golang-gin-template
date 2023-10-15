package routes

import (
	"myapp/controllers"

	"github.com/gin-gonic/gin"
)

func AlbumRoutes(router gin.Engine) {
	album := router.Group("/albums")

	album.GET("/", controllers.GetAlbums)
	album.POST("/", controllers.PostAlbums)
}
