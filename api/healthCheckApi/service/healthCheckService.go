package service

import (
	"github.com/gin-gonic/gin"
	"github.com/rzabhd80/healthCheck/api/healthCheckApi/handlers"
	"github.com/rzabhd80/healthCheck/api/healthCheckApi/repository"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, database *gorm.DB) {
	apiRepository := repository.NewAPIRepository(database)
	apiHandler := handlers.NewAPIHandler(apiRepository)
	apiRoutes := router.Group("/api")
	{
		apiRoutes.POST("/", apiHandler.CreateAPI)
		apiRoutes.GET("/", apiHandler.GetAPIs)
		apiRoutes.DELETE("/:id", apiHandler.DeleteAPI)
	}
}
