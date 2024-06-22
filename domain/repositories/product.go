package repositories

import "saas-billing/domain/entities"

type ProductRepository interface {
	FindByID(id string) (*entities.Product, error)
	Create(product *entities.Product) error
}
