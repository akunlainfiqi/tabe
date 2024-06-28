package repositories

import "saas-billing/domain/entities"

type BillsRepository interface {
	GetByID(id string) (*entities.Bills, error)
	GetUnpaidBillsAfterDueDate() ([]*entities.Bills, error)
	Create(billing *entities.Bills) error
	Update(billing *entities.Bills) error
}
