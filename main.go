package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rzabhd80/healthCheck/api/healthCheckApi/service"
	"github.com/rzabhd80/healthCheck/internals"
	"log"
)

func main() {
	config, _ := internals.LoadConfig()
	log.Printf("port%s", config.ServerPort)
	engine := gin.Default()
	db := internals.Setup(config)
	internals.Migrate(db)
	service.SetupApp(engine, db)
	err := engine.Run(":" + config.ServerPort)
	if err != nil {
		return
	}
}
