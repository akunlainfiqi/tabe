package entities

import (
	"saas-billing/errors"
	"time"
)

type BillType string
type BillStatus string

const (
	BillStatusWaitingPayment BillStatus = "waiting_payment"
	BillStatusPaid           BillStatus = "paid"
	BillStatusOverdue        BillStatus = "overdue"
	BillStatusCancelled      BillStatus = "cancelled"
)

func mapBillStatus(status string) (BillStatus, error) {
	switch status {
	case "waiting_payment":
		return BillStatusWaitingPayment, nil
	case "paid":
		return BillStatusPaid, nil
	case "overdue":
		return BillStatusOverdue, nil
	case "cancelled":
		return BillStatusCancelled, nil
	default:
		return "", errors.ErrInvalidBillStatus
	}
}

type Bills struct {
	id             string
	organizationId string
	tenantId       string
	priceId        string
	status         BillStatus
	amount         int64
	balanceUsed    int64
	dueDate        int64
	billType       BillType
}

const (
	BillTypeNewSubscription BillType = "new_subscription"
	BillTypeRenewal         BillType = "renewal"
	BillTypeUpgrade         BillType = "upgrade"
	BillTypeDowngrade       BillType = "downgrade"
	BillTypeAddBalance      BillType = "add_balance"
)

func mapBillType(billType string) (BillType, error) {
	switch billType {
	case "new_subscription":
		return BillTypeNewSubscription, nil
	case "renewal":
		return BillTypeRenewal, nil
	case "upgrade":
		return BillTypeUpgrade, nil
	case "downgrade":
		return BillTypeDowngrade, nil
	case "add_balance":
		return BillTypeAddBalance, nil
	default:
		return "", errors.ErrInvalidBillType
	}
}

func BuildBills(
	id,
	organizationId,
	tenantId,
	priceId string,
	status string,
	amount,
	balanceUsed,
	dueDate int64,
	billType string,
) *Bills {
	BillStatus, _ := mapBillStatus(status)
	billTypes, _ := mapBillType(billType)
	return &Bills{
		id:             id,
		organizationId: organizationId,
		tenantId:       tenantId,
		priceId:        priceId,
		status:         BillStatus,
		amount:         amount,
		balanceUsed:    balanceUsed,
		dueDate:        dueDate,
		billType:       billTypes,
	}
}

func NewBills(
	id,
	organizationId,
	tenantId,
	priceId string,
	amount,
	balanceUsed,
	dueDate int64,
	billType string,
) (*Bills, error) {
	var billStatus BillStatus
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

	billTypes, err := mapBillType(billType)
	if err != nil {
		return nil, err
	}

	return &Bills{
		id:             id,
		organizationId: organizationId,
		tenantId:       tenantId,
		priceId:        priceId,
		status:         billStatus,
		dueDate:        dueDate,
		amount:         amount,
		balanceUsed:    balanceUsed,
		billType:       billTypes,
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

func (b *Bills) Status() BillStatus {
	return b.status
}

func (b *Bills) SetStatus(status BillStatus) {
	b.status = status
}

func (b *Bills) Pay(amount int64) error {
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

func (b *Bills) Settle() {
	b.status = BillStatusPaid
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

func (b *Bills) Total() int64 {
	return b.amount - b.balanceUsed
}
func (b *Bills) Amount() int64 {
	return b.amount
}

func (b *Bills) BalanceUsed() int64 {
	return b.balanceUsed
}

func (b *Bills) DueDate() int64 {
	return b.dueDate
}

func (b *Bills) BillType() BillType {
	return b.billType
}

func (b *Bills) PriceID() string {
	return b.priceId
}
