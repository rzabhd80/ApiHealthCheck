package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rzabhd80/healthCheck/api/healthCheckApi/service"
	"github.com/rzabhd80/healthCheck/internals"
)

func main() {
	config, _ := internals.LoadConfig()
	engine := gin.Default()
	db := internals.Setup(config)
	internals.Migrate(db)
	service.SetupRoutes(engine, db)
	err := engine.Run(":" + config.ServerPort)
	if err != nil {
		return
	}
}
