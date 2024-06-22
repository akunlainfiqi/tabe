package repositories

import "saas-billing/domain/entities"

type TransactionRepository interface {
	GetByID(id string) (*entities.Transaction, error)
	Create(transaction *entities.Transaction) error
	Update(transaction *entities.Transaction) error
}
