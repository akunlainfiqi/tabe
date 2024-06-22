package pgsql

import (
	"saas-billing/domain/entities"
	"saas-billing/domain/repositories"
	"saas-billing/errors"
	"time"

	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) repositories.TransactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) GetByID(id string) (*entities.Transaction, error) {
	var dto struct {
		ID                   string
		BillsID              string
		OrganizationId       string
		Amount               int64
		TransactionType      string
		TransactionTimestamp int64
	}
	r.db.Raw(`
		SELECT 
			t.id,
			t.bills_id,
			t.organization_id,
			t.amount,
			t.transaction_type,
			t.transaction_timestamp
		FROM transactions t
		WHERE t.id = ?
		`, id).Scan(&dto)

	transaction := entities.BuildTransaction(dto.ID, dto.BillsID, dto.OrganizationId, dto.Amount, dto.TransactionType, dto.TransactionTimestamp)

	if transaction.ID() == "" {
		return nil, errors.ErrTransactionNotFound
	}
	return transaction, nil
}

func (r *transactionRepository) Create(transaction *entities.Transaction) error {
	if err := r.db.Exec(`
		INSERT INTO transactions (id, bills_id, organization_id, amount, transaction_type, transaction_timestamp, created_at, updated_at)
		VALUES (@id, @bills_id, @organization_id, @amount, @transaction_type, @transaction_timestamp, @date, @date)
	`, map[string]interface{}{
		"id":                    transaction.ID(),
		"bills_id":              transaction.BillsID(),
		"organization_id":       transaction.OrganizationID(),
		"amount":                transaction.Amount(),
		"transaction_type":      transaction.TransactionType(),
		"transaction_timestamp": transaction.TransactionTimestamp(),
		"date":                  time.Now().Unix(),
	}).Error; err != nil {
		return err
	}
	return nil
}

func (r *transactionRepository) Update(transaction *entities.Transaction) error {
	if err := r.db.Exec(`
		UPDATE transactions
		SET
			bills_id = @bills_id,
			organization_id = @organization_id,
			amount = @amount,
			transaction_type = @transaction_type,
			transaction_timestamp = @transaction_timestamp,
			updated_at = @date
		WHERE id = @id
	`, map[string]interface{}{
		"id":                    transaction.ID(),
		"bills_id":              transaction.BillsID(),
		"organization_id":       transaction.OrganizationID(),
		"amount":                transaction.Amount(),
		"transaction_type":      transaction.TransactionType(),
		"transaction_timestamp": transaction.TransactionTimestamp(),
		"date":                  time.Now().Unix(),
	}).Error; err != nil {
		return err
	}
	return nil
}
