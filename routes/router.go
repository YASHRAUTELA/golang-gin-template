package routes

import (
	"os"
	"server/middleware"
	"strings"

	"github.com/gin-gonic/gin"
)

func InitRoute() *gin.Engine {
	environment := os.Getenv("DEBUG")
	if environment == "true" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	allowedHosts := os.Getenv("ALLOWED_HOSTS")
	router := gin.New()
	router.SetTrustedProxies(strings.Split(allowedHosts, ","))
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.InitCORSMiddleware())
	DefaultRoutes(router)
	AuthRoutes(router)
	UserRoutes(router)
	WorkspaceRoutes(router)
	SessionRoutes(router)
	return router
}
