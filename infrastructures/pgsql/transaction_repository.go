package pgsql

import (
	"saas-billing/domain/entities"
	"saas-billing/domain/repositories"
	"saas-billing/errors"

	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) repositories.TransactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) GetByID(id string) (*entities.Transaction, error) {
	var transaction entities.Transaction
	if err := r.db.Where("id = ?", id).First(&transaction).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrTransactionNotFound
		}
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) Create(transaction *entities.Transaction) error {
	if err := r.db.Create(transaction).Error; err != nil {
		return err
	}
	return nil
}

func (r *transactionRepository) Update(transaction *entities.Transaction) error {
	if err := r.db.Save(transaction).Error; err != nil {
		return err
	}
	return nil
}
