package pgsql

import (
	"saas-billing/domain/entities"
	"saas-billing/domain/repositories"

	"gorm.io/gorm"
)

// ProductRepository is a repository for managing Product entities
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new ProductRepository
func NewProductRepository(db *gorm.DB) repositories.ProductRepository {
	return &ProductRepository{db}
}

// FindByID finds a product by its ID
func (pr *ProductRepository) FindByID(id string) (*entities.Product, error) {
	var product entities.Product
	if err := pr.db.Where("id = ?", id).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

// Create creates a new product
func (pr *ProductRepository) Create(product *entities.Product) error {
	return pr.db.Create(product).Error
}
