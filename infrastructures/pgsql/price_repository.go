package pgsql

import (
	"saas-billing/domain/entities"
	"saas-billing/domain/repositories"

	"gorm.io/gorm"
)

type PriceRepository struct {
	db *gorm.DB
}

func NewPriceRepository(db *gorm.DB) repositories.PriceRepository {
	return &PriceRepository{db}
}

// FindByID finds a price by its ID
func (pr *PriceRepository) FindByID(id string) (*entities.Price, error) {
	var price entities.Price
	if err := pr.db.Where("id = ?", id).First(&price).Error; err != nil {
		return nil, err
	}
	return &price, nil
}

// Create creates a new price
func (pr *PriceRepository) Create(price *entities.Price) error {
	return pr.db.Create(price).Error
}
