package pgsql

import (
	"saas-billing/domain/entities"
	"saas-billing/domain/repositories"
	"saas-billing/errors"

	"gorm.io/gorm"
)

type billsRepository struct {
	db *gorm.DB
}

func NewBillsRepository(db *gorm.DB) repositories.BillsRepository {
	return &billsRepository{db}
}

func (r *billsRepository) GetByID(id string) (*entities.Bills, error) {
	var bills entities.Bills
	if err := r.db.Where("id = ?", id).First(&bills).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrBillsNotFound
		}
		return nil, err
	}
	return &bills, nil
}

func (r *billsRepository) Create(billing *entities.Bills) error {
	if err := r.db.Create(billing).Error; err != nil {
		return err
	}
	return nil
}

func (r *billsRepository) Update(billing *entities.Bills) error {
	if err := r.db.Save(billing).Error; err != nil {
		return err
	}
	return nil
}
