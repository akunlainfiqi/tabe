package pgsql

import (
	"saas-billing/app/queries"
	"saas-billing/errors"

	"gorm.io/gorm"
)

type billQuery struct {
	db *gorm.DB
}

func NewBillQuery(db *gorm.DB) queries.BillQuery {
	return &billQuery{db}
}

func (q *billQuery) GetByID(id string) (queries.Bill, error) {
	var bill queries.Bill

	if err := q.db.Raw(`
		SELECT
			b.id,
			b.organization_id,
			b.tenant_id,
			t.name AS tenant_name,
			b.status,
			b.due_date,
			b.amount,
			b.balance_used,
			b.created_at
		FROM
			bills b
		JOIN
			tenants t
		ON
			b.tenant_id = t.id
		WHERE
			b.id = ?
	`, id).Scan(&bill).Error; err != nil {
		return queries.Bill{}, err
	}

	if bill.Id == "" {
		return queries.Bill{}, errors.ErrBillsNotFound
	}

	return bill, nil
}

func (q *billQuery) GetByOrganizationID(organizationID string) ([]queries.Bill, error) {
	var bills []queries.Bill

	if err := q.db.Raw(`
		SELECT
			b.id,
			b.organization_id,
			b.tenant_id,
			t.name AS tenant_name,
			b.status,
			b.due_date,
			b.amount,
			b.balance_used,
			b.created_at
		FROM
			bills b
		JOIN
			tenants t
		ON
			b.tenant_id = t.id
		WHERE
			b.organization_id = ?
	`, organizationID).Scan(&bills).Error; err != nil {
		return []queries.Bill{}, err
	}

	return bills, nil
}
