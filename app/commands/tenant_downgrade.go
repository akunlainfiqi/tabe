package commands

import (
	"encoding/json"
	"saas-billing/app/services"
	"saas-billing/domain/entities"
	"saas-billing/domain/repositories"
	"saas-billing/errors"
	"time"

	"github.com/google/uuid"
)

type TenantDowngradeRequest struct {
	TenantID string
	PriceID  string
}

type TenantDowngradeCommand struct {
	tenantRepository repositories.TenantRepository
	priceRepository  repositories.PriceRepository
	orgRepository    repositories.OrganizationRepository
	billsRepository  repositories.BillsRepository
	midtransService  services.Midtrans
	publisherService services.Publisher
}

func NewTenantDowngradeCommand(
	tenantRepository repositories.TenantRepository,
	priceRepository repositories.PriceRepository,
	orgRepository repositories.OrganizationRepository,
	billsRepository repositories.BillsRepository,
	midtransService services.Midtrans,
	publisherService services.Publisher,
) *TenantDowngradeCommand {
	return &TenantDowngradeCommand{
		tenantRepository: tenantRepository,
		priceRepository:  priceRepository,
		orgRepository:    orgRepository,
		billsRepository:  billsRepository,
		midtransService:  midtransService,
		publisherService: publisherService,
	}
}

func (c *TenantDowngradeCommand) Execute(req *TenantDowngradeRequest) (interface{}, error) {
	tenant, err := c.tenantRepository.GetByID(req.TenantID)
	if err != nil {
		return nil, err
	}

	if tenant.PriceID() == req.PriceID {
		return nil, errors.ErrInvalidSameProductPrice
	}

	org, err := c.orgRepository.GetByID(tenant.OrganizationID())
	if err != nil {
		return nil, err
	}

	newPrice, err := c.priceRepository.GetByID(req.PriceID)
	if err != nil {
		return nil, err
	}
	oldPrice, err := c.priceRepository.GetByID(tenant.PriceID())
	if err != nil {
		return nil, err
	}

	// Check if the new price is lower than the old price
	if newPrice.Price() > oldPrice.Price() {
		return nil, errors.ErrInvalidPrice
	}

	if newPrice.Product().TierIndex() > oldPrice.Product().TierIndex() {
		return nil, errors.ErrInvalidTierIndex
	}

	var balanceFairUsageRefund int64
	if oldPrice.Recurrence() == entities.ProductRecurrenceMonthly {
		//parse activeuntil from unix to time.Time
		activeUntil := time.Unix(tenant.ActiveUntil(), 0)
		//calculate the balance fair usage refund
		// activeuntil - now = remaining days * old price / 30
		activeUntilHours := time.Until(activeUntil).Hours()
		balanceFairUsageRefund = int64(int64(activeUntilHours) * oldPrice.Price() / 30)
	} else {
		//parse activeuntil from unix to time.Time
		activeUntil := time.Unix(tenant.ActiveUntil(), 0)
		//calculate the balance fair usage refund
		// activeuntil - now = remaining days * old price / 365
		activeUntilHours := time.Until(activeUntil).Hours()
		balanceFairUsageRefund = int64(int64(activeUntilHours) * oldPrice.Price() / 365)
	}

	// Update the organization balance
	org.SetBalance(org.Balance() + balanceFairUsageRefund)

	// Create a new bill for the balance fair usage refund
	billId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	if newPrice.Price() < org.Balance() {
		balanceUsed := newPrice.Price()
		bill, err := entities.NewBills(
			billId.String(),
			org.ID(),
			tenant.ID(),
			newPrice.ID(),
			newPrice.Price(),
			balanceUsed,
			time.Now().AddDate(0, 0, 10).Unix(),
			string(entities.BillTypeDowngrade),
		)
		if err != nil {
			return nil, err
		}

		if err := c.billsRepository.Create(bill); err != nil {
			return nil, err
		}

		org.SetBalance(org.Balance() - balanceUsed)

		if err := c.orgRepository.Update(org); err != nil {
			return nil, err
		}

		tenant.SetPriceID(newPrice.ID())

		if newPrice.Recurrence() == entities.ProductRecurrenceMonthly {
			tenant.SetActiveUntil(time.Now().AddDate(0, 1, 0).Unix())
		} else if newPrice.Recurrence() == entities.ProductRecurrenceYearly {
			tenant.SetActiveUntil(time.Now().AddDate(1, 0, 0).Unix())
		}

		if err := c.tenantRepository.Update(tenant); err != nil {
			return nil, err
		}

		pl := services.TenantPaidPayload{
			TenantID:  bill.TenantID(),
			ProductID: newPrice.Product().ID(),
			Timestamp: time.Now(),
		}

		payload, err := json.Marshal(pl)
		if err != nil {
			return nil, err
		}

		if err := c.publisherService.Publish("billing_paid", payload); err != nil {
			return nil, err
		}

		return nil, nil
	}

	//if the balance is not enough
	balanceUsed := org.Balance()
	bill, err := entities.NewBills(
		billId.String(),
		org.ID(),
		tenant.ID(),
		newPrice.ID(),
		newPrice.Price(),
		balanceUsed,
		time.Now().AddDate(0, 0, 10).Unix(),
		string(entities.BillTypeDowngrade),
	)

	if err != nil {
		return nil, err
	}

	if err := c.billsRepository.Create(bill); err != nil {
		return nil, err
	}

	org.SetBalance(0)

	if err := c.orgRepository.Update(org); err != nil {
		return nil, err
	}

	res, err := c.midtransService.CreateTransaction(bill.ID(), bill.Total())
	if err != nil {
		return nil, err
	}

	return res, nil
}
