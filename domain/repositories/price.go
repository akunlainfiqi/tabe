package repositories

import "saas-billing/domain/entities"

type PriceRepository interface {
	FindByID(id string) (*entities.Price, error)
	Create(price *entities.Price) error
}
