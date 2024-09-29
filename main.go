package main

import (
	"fmt"
	"os"
	"server/config"
	"server/routes"
)

func main() {
	config.EnvLoad()
	config.InitDBConnection()

	// @title Gin Postgres Swagger Example API
	// @version 1.0
	// @description This is a sample server.
	// @termsOfService http://swagger.io/terms/

	// @contact.name API Support
	// @contact.url http://www.swagger.io/support
	// @contact.email support@swagger.io

	// @license.name Apache 2.0
	// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

	// @host localhost:9000
	// @BasePath /

	// @schemes http

	// @securityDefinitions.apikey BearerAuth
	// @in header
	// @name Authorization
	// @description "Type 'Bearer' followed by a space and your JWT token."

	// @Summary Ping example
	// @Description do ping
	// @ID ping
	// @Accept  json
	// @Produce  json
	// @Success 200 {object} map[string]string
	// @Router /ping [get]

	r := routes.InitRoute()
	r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}
