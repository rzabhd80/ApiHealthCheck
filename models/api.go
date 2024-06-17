package models

import (
	"gorm.io/gorm"
	"time"
)

type API struct {
	gorm.Model
	ID                  uint       `gorm:"primaryKey"`
	HealthCheckInterval int        `json:"interval"`
	RequestURL          string     `json:"requestURL"`
	RequestMethod       string     `json:"requestMethod"`
	RequestHeaders      *string    `json:"requestHeaders"`
	RequestBody         *string    `json:"requestBody"`
	LastStatus          *int       `json:"lastStatus"`
	ShouldBeChecked     bool       `json:"shouldBeChecked"`
	LastChecked         *time.Time `json:"lastChecked"`
}
