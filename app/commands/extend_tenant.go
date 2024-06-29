package commands

import (
	"saas-billing/app/services"
	"saas-billing/domain/entities"
	"saas-billing/domain/repositories"
	"time"

	"github.com/google/uuid"
)

type ExtendTenantCommandRequest struct {
	TenantID string
	UserID   string
}

type ExtendTenantCommand struct {
	tenantRepository repositories.TenantRepository
	orgRepository    repositories.OrganizationRepository
	priceRepository  repositories.PriceRepository
	billsRepository  repositories.BillsRepository

	midtransService services.Midtrans
}

func NewExtendTenantCommand(
	tenantRepository repositories.TenantRepository,
	orgRepository repositories.OrganizationRepository,
	priceRepository repositories.PriceRepository,
	billsRepository repositories.BillsRepository,
	midtransService services.Midtrans,
) *ExtendTenantCommand {
	return &ExtendTenantCommand{
		tenantRepository: tenantRepository,
		orgRepository:    orgRepository,
		priceRepository:  priceRepository,
		billsRepository:  billsRepository,
		midtransService:  midtransService,
	}
}

func (c *ExtendTenantCommand) Do(req *ExtendTenantCommandRequest) (interface{}, error) {
	// Get tenant by ID
	tenant, err := c.tenantRepository.GetByID(req.TenantID)
	if err != nil {
		return nil, err
	}

	// Get organization by ID
	org, err := c.orgRepository.GetByID(tenant.OrganizationID())
	if err != nil {
		return nil, err
	}

	// Get price by ID
	price, err := c.priceRepository.GetByID(tenant.PriceID())
	if err != nil {
		return nil, err
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	balanceUsed := int64(0)

	if org.Balance() > 0 {
		if org.Balance() >= price.Price() {
			balanceUsed = price.Price()
		} else {
			balanceUsed = org.Balance()
		}
	}

	org.SetBalance(org.Balance() - balanceUsed)

	// Create bills
	billing, err := entities.NewBills(
		id.String(),
		org.ID(),
		tenant.ID(),
		price.Price(),
		balanceUsed,
		time.Now().AddDate(0, 0, 1).Unix(),
		entities.BillTypeRenewal,
	)
	if err != nil {
		return nil, err
	}

	if err := c.billsRepository.Create(billing); err != nil {
		return nil, err
	}

	if balanceUsed > 0 {
		if err := c.orgRepository.Update(org); err != nil {
			return nil, err
		}
	}

	res, err := c.midtransService.CreateTransaction(billing.ID(), billing.Total())
	if err != nil {
		return nil, err
	}

	return res, nil
}
