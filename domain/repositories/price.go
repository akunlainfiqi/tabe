package repositories

import "saas-billing/domain/entities"

type PriceRepository interface {
	GetByID(id string) (*entities.Price, error)
	Create(price *entities.Price) error
}
