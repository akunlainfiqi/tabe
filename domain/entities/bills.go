package entities

import (
	"saas-billing/errors"
	"time"
)

type Bills struct {
	id             string 
	organizationId string
	tenantId       string
	status         string
	amount         float64
	balanceUsed    float64
	dueDate        int64
}

const (
	BillStatusWaitingPayment = "waiting_payment"
	BillStatusPaid           = "paid"
	BillStatusOverdue        = "overdue"
	BillStatusCancelled      = "cancelled"
)

func NewBills(
	id,
	organizationId,
	tenantId string,
	amount,
	balanceUsed float64,
	dueDate int64,
) (*Bills, error) {
	var billStatus string
	if balanceUsed == 0 {
		billStatus = BillStatusWaitingPayment
	}
	if amount > balanceUsed {
		billStatus = BillStatusWaitingPayment
	}
	if balanceUsed > 0 && balanceUsed == amount {
		billStatus = BillStatusPaid
	}
	if balanceUsed > 0 && balanceUsed > amount {
		return nil, errors.ErrInvalidBillAmount
	}
	return &Bills{
		id:             id,
		organizationId: organizationId,
		tenantId:       tenantId,
		amount:         amount,
		status:         billStatus,
		balanceUsed:    balanceUsed,
		dueDate:        dueDate,
	}, nil
}

func (b *Bills) ID() string {
	return b.id
}

func (b *Bills) OrganizationID() string {
	return b.organizationId
}

func (b *Bills) TenantID() string {
	return b.tenantId
}

func (b *Bills) Status() string {
	return b.status
}

func (b *Bills) SetStatus(status string) {
	b.status = status
}

func (b *Bills) Pay(amount float64) error {
	if amount > b.amount {
		return errors.ErrInvalidBillAmount
	}
	if amount > b.amount-b.balanceUsed {
		return errors.ErrInvalidBillAmount
	}
	b.balanceUsed += amount
	if b.balanceUsed == b.amount {
		b.status = BillStatusPaid
	}
	return nil
}

func (b *Bills) Cancel() {
	b.status = BillStatusCancelled
}

func (b *Bills) Expire() {
	b.status = BillStatusOverdue
}

func (b *Bills) IsOverdue() bool {
	return b.dueDate < time.Now().Unix()
}

func (b *Bills) Amount() float64 {
	return b.amount
}

func (b *Bills) BalanceUsed() float64 {
	return b.balanceUsed
}

func (b *Bills) DueDate() int64 {
	return b.dueDate
}
