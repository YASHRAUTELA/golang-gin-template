package config

import (
	_ "server/docs" // This is necessary for the swagger docs

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SwaggerConfig() gin.HandlerFunc {
	return ginSwagger.WrapHandler(swaggerFiles.Handler)
}
