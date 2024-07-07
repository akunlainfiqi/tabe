package pgsql

import (
	"saas-billing/app/queries"

	"gorm.io/gorm"
)

type transactionQuery struct {
	db *gorm.DB
}

func NewTransactionQuery(db *gorm.DB) queries.TransactionQuery {
	return &transactionQuery{db}
}

func (q *transactionQuery) FindByOrgID(id string) ([]queries.Transaction, error) {
	var temp []queries.Transaction
	if err := q.db.Raw(`
		SELECT
			t.id,
			t.organization_id,
			t.amount,
			t.transaction_type,
			t.transaction_timestamp
		FROM
			transactions t
		WHERE
			t.organization_id = ?
	`, id).Scan(&temp).Error; err != nil {
		return []queries.Transaction{}, err
	}

	return temp, nil
}

func (q *transactionQuery) FindByBillsID(id string) (queries.Transaction, error) {
	var temp queries.Transaction
	if err := q.db.Raw(`
		SELECT
			t.id,
			t.bills_id,
			t.organization_id,
			t.amount,
			t.transaction_type,
			t.transaction_timestamp
		FROM
			transactions t
		WHERE
			t.bills_id = ?
	`, id).Scan(&temp).Error; err != nil {
		return queries.Transaction{}, err
	}

	return temp, nil
}
