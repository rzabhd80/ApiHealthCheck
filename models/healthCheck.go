package models

import (
	"gorm.io/gorm"
)

type HealthCheck struct {
	gorm.Model
	URL    string `json:"url"`
	ApiId  uint   `json:"api_id"`
	Status string `json:"status"`
}
