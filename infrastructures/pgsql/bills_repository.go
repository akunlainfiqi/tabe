package pgsql

import (
	"saas-billing/domain/entities"
	"saas-billing/domain/repositories"
	"saas-billing/errors"
	"time"

	"gorm.io/gorm"
)

type billsRepository struct {
	db *gorm.DB
}

func NewBillsRepository(db *gorm.DB) repositories.BillsRepository {
	return &billsRepository{db}
}

func (r *billsRepository) GetByID(id string) (*entities.Bills, error) {
	var dto struct {
		ID             string
		OrganizationID string
		TenantID       string
		Status         string
		DueDate        int64
		Amount         int64
		BalanceUsed    int64
		BillType       string
	}
	r.db.Raw(`
		SELECT
			b.id,
			b.organization_id,
			b.tenant_id,
			b.status,
			b.due_date,
			b.amount,
			b.balance_used,
			b.bill_type
		FROM bills b
		WHERE b.id = ?
		`, id).Scan(&dto)

	billing := entities.BuildBills(dto.ID, dto.OrganizationID, dto.TenantID, dto.Status, dto.DueDate, dto.Amount, dto.BalanceUsed, dto.BillType)

	if billing.ID() == "" {
		return nil, errors.ErrBillsNotFound
	}

	return billing, nil

}

func (r *billsRepository) Create(billing *entities.Bills) error {
	err := r.db.Exec(`
		INSERT INTO bills (id, organization_id, tenant_id, status, due_date, amount, balance_used, bill_type, created_at, updated_at)
		VALUES (@id, @organization_id, @tenant_id, @status, @due_date, @amount, @balance_used, @bill_type, @now, @now)
	`, map[string]interface{}{
		"id":              billing.ID(),
		"organization_id": billing.OrganizationID(),
		"tenant_id":       billing.TenantID(),
		"status":          billing.Status(),
		"due_date":        billing.DueDate(),
		"amount":          billing.Amount(),
		"balance_used":    billing.BalanceUsed(),
		"bill_type":       billing.BillType(),
		"now":             time.Now().Unix(),
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *billsRepository) Update(billing *entities.Bills) error {
	err := r.db.Exec(`
		UPDATE bills
		SET
			organization_id = @organization_id,
			tenant_id = @tenant_id,
			status = @status,
			due_date = @due_date,
			amount = @amount,
			balance_used = @balance_used,
			bill_type = @bill_type,
			updated_at = @now
		WHERE id = @id
	`, map[string]interface{}{
		"id":              billing.ID(),
		"organization_id": billing.OrganizationID(),
		"tenant_id":       billing.TenantID(),
		"status":          billing.Status(),
		"due_date":        billing.DueDate(),
		"amount":          billing.Amount(),
		"balance_used":    billing.BalanceUsed(),
		"bill_type":       billing.BillType(),
		"now":             time.Now().Unix(),
	}).Error
	return err
}

func (r *billsRepository) GetUnpaidBillsAfterDueDate() ([]*entities.Bills, error) {
	var dtos []struct {
		ID             string
		OrganizationID string
		TenantID       string
		Status         string
		DueDate        int64
		Amount         int64
		BalanceUsed    int64
		BillType       string
	}
	r.db.Raw(`
		SELECT
			b.id,
			b.organization_id,
			b.tenant_id,
			b.status,
			b.due_date,
			b.amount,
			b.balance_used,
			b.bill_type
		FROM bills b
		WHERE b.status = 'waiting_payment' AND b.due_date < ?
		`, time.Now().Unix()).Scan(&dtos)

	billings := make([]*entities.Bills, 0)
	for _, dto := range dtos {
		billings = append(billings, entities.BuildBills(
			dto.ID,
			dto.OrganizationID,
			dto.TenantID,
			dto.Status,
			dto.Amount,
			dto.BalanceUsed,
			dto.DueDate,
			dto.BillType))
	}

	return billings, nil
}
