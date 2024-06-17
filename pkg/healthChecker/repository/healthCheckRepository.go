package repository

import (
	"github.com/rzabhd80/healthCheck/models"
	"gorm.io/gorm"
	"log"
)

type HealthCheckRepository interface {
	LogHealthCheck(apiID uint, url string, status string) error
}

type healthCheckRepository struct {
	db *gorm.DB
}

func NewHealthCheckRepository(db *gorm.DB) HealthCheckRepository {
	return &healthCheckRepository{db}
}

func (r *healthCheckRepository) LogHealthCheck(apiID uint, url string, status string) error {
	healthCheck := models.HealthCheck{
		ApiID:  apiID,
		URL:    url,
		Status: status,
	}
	
	if err := r.db.Create(&healthCheck).Error; err != nil {
		log.Printf("Failed to log health check: %v", err)
		return err
	}
	return nil
}
