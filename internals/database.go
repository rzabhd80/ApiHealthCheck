package internals

import (
	"fmt"
	"github.com/rzabhd80/healthCheck/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func Setup(cfg *Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	return db
}

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(&models.API{}, &models.HealthCheck{})
	if err != nil {
		return
	}
}
