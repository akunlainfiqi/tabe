package pgsql

import (
	"saas-billing/domain/entities"
	"saas-billing/domain/repositories"

	"gorm.io/gorm"
)

type AppsRepository struct {
	db *gorm.DB
}

func NewAppsRepository(db *gorm.DB) repositories.AppsRepository {
	return &AppsRepository{db}
}

// FindByID finds an app by its ID
func (ar *AppsRepository) FindByID(id string) (*entities.Apps, error) {
	var app entities.Apps
	if err := ar.db.Where("id = ?", id).First(&app).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

// Create creates a new app
func (ar *AppsRepository) Create(app *entities.Apps) error {
	return ar.db.Create(app).Error
}
