package service

import (
	"github.com/gin-gonic/gin"
	"github.com/rzabhd80/healthCheck/api/healthCheckApi/handlers"
	"github.com/rzabhd80/healthCheck/api/healthCheckApi/repository"
	"github.com/rzabhd80/healthCheck/pkg/healthChecker"
	healthCheckRepo "github.com/rzabhd80/healthCheck/pkg/healthChecker/repository"
	"gorm.io/gorm"
	"os"
	"time"
)

func SetupRoutes(router *gin.Engine, database *gorm.DB) {
	apiRepository := repository.NewAPIRepository(database)
	apiHandler := handlers.NewAPIHandler(apiRepository)
	healthCheckRepository := healthCheckRepo.NewHealthCheckRepository(database)
	healthCheck := healthChecker.NewHealthChecker(apiRepository, os.Getenv("SLACK_WEBHOOK_URL"),
		1*time.Minute, healthCheckRepository)
	healthCheck.Start(5)
	apiRoutes := router.Group("/api")
	{
		apiRoutes.POST("/", apiHandler.CreateAPI)
		apiRoutes.GET("/", apiHandler.GetAPIs)
		apiRoutes.DELETE("/:id", apiHandler.DeleteAPI)
	}
}
