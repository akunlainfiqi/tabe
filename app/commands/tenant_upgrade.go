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

type TenantUpgradeRequest struct {
	TenantID string
	PriceID  string
}

type TenantUpgradeCommand struct {
	tenantRepository       repositories.TenantRepository
	priceRepository        repositories.PriceRepository
	billRepository         repositories.BillsRepository
	organizationRepository repositories.OrganizationRepository
	midtransService        services.Midtrans
	publisherService       services.Publisher
}

func NewTenantUpgradeCommand(
	tenantRepository repositories.TenantRepository,
	priceRepository repositories.PriceRepository,
	billRepository repositories.BillsRepository,
	organizationRepository repositories.OrganizationRepository,
	midtransService services.Midtrans,
	publisherService services.Publisher,
) *TenantUpgradeCommand {
	return &TenantUpgradeCommand{
		tenantRepository:       tenantRepository,
		priceRepository:        priceRepository,
		billRepository:         billRepository,
		organizationRepository: organizationRepository,
		midtransService:        midtransService,
		publisherService:       publisherService,
	}
}

func (c *TenantUpgradeCommand) Execute(req *TenantUpgradeRequest) (interface{}, error) {
	tenant, err := c.tenantRepository.GetByID(req.TenantID)
	if err != nil {
		return nil, err
	}

	org, err := c.organizationRepository.GetByID(tenant.OrganizationID())
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

	if tenant.PriceID() == req.PriceID {
		return nil, errors.ErrInvalidSameProductPrice
	}

	if newPrice.Product().TierIndex() < oldPrice.Product().TierIndex() {
		return nil, errors.ErrInvalidTierIndex
	}

	if newPrice.Product().App().ID() != oldPrice.Product().App().ID() {
		return nil, errors.ErrMismatchedApp
	}

	//assuming that the new price is always higher than the old price
	if newPrice.Price() < oldPrice.Price() {
		return nil, errors.ErrInvalidPrice
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	totalPrice := newPrice.Price() - oldPrice.Price()
	
	if org.Balance() > totalPrice {
		balanceUsed := totalPrice

		bill, err := entities.NewBills(
			id.String(),
			org.ID(),
			tenant.ID(),
			newPrice.ID(),
			newPrice.Price(),
			balanceUsed,
			time.Now().Add(time.Hour).Unix(),
			string(entities.BillTypeUpgrade),
		)

		if err != nil {
			return nil, err
		}

		if err := c.billRepository.Create(bill); err != nil {
			return nil, err
		}

		org.SetBalance(org.Balance() - balanceUsed)

		if err := c.organizationRepository.Update(org); err != nil {
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

	org.SetBalance(0)

	if err := c.organizationRepository.Update(org); err != nil {
		return nil, err
	}

	bill, err := entities.NewBills(
		id.String(),
		org.ID(),
		tenant.ID(),
		newPrice.ID(),
		newPrice.Price(),
		balanceUsed,
		time.Now().Add(time.Hour).Unix(),
		string(entities.BillTypeUpgrade),
	)

	if err != nil {
		return nil, err
	}

	if err := c.billRepository.Create(bill); err != nil {
		return nil, err
	}

	res, err := c.midtransService.CreateTransaction(bill.ID(), bill.Total())
	if err != nil {
		return nil, err
	}
	return res, nil
}
