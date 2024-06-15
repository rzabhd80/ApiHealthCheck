package repository

import (
	"github.com/rzabhd80/healthCheck/models"
	"gorm.io/gorm"
)

type APIRepository interface {
	Create(api *models.API) error
	GetAll() ([]models.API, error)
	GetByID(id uint) (*models.API, error)
	GetAllActive() ([]models.API, error)
	Delete(id uint) error
	Update(api *models.API) error
}

type ApiRepository struct {
	db *gorm.DB
}

func NewAPIRepository(db *gorm.DB) *ApiRepository {
	return &ApiRepository{db: db}
}

func (r *ApiRepository) Create(api *models.API) error {
	return r.db.Create(api).Error
}

func (r *ApiRepository) GetAll() ([]models.API, error) {
	var apis []models.API
	err := r.db.Find(&apis).Error
	return apis, err
}

func (r *ApiRepository) GetByID(id uint) (*models.API, error) {
	var api models.API
	err := r.db.First(&api, id).Error
	return &api, err
}

func (r *ApiRepository) Delete(id uint) error {
	return r.db.Delete(&models.API{}, id).Error
}

func (r *ApiRepository) Update(api *models.API) error {
	return r.db.Save(api).Error
}

func (r *ApiRepository) GetAllActive() ([]models.API, error) {
	var apis []models.API
	err := r.db.Where("should_be_checked = ?", true).Find(&apis).Error
	return apis, err
}
