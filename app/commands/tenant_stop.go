package commands

import (
	"saas-billing/domain/entities"
	"saas-billing/domain/repositories"
	"time"
)

type TenantStopRequest struct {
	TenantID string
}

type TenantStopCommand struct {
	tenantRepository repositories.TenantRepository
	orgRepository    repositories.OrganizationRepository
	billsRepository  repositories.BillsRepository
	priceRepository  repositories.PriceRepository
}

func NewTenantStopCommand(
	tenantRepository repositories.TenantRepository,
	orgRepository repositories.OrganizationRepository,
	billsRepository repositories.BillsRepository,
	priceRepository repositories.PriceRepository,
) *TenantStopCommand {
	return &TenantStopCommand{
		tenantRepository: tenantRepository,
		orgRepository:    orgRepository,
		billsRepository:  billsRepository,
		priceRepository:  priceRepository,
	}
}

func (c *TenantStopCommand) Execute(req *TenantStopRequest) error {
	tenant, err := c.tenantRepository.GetByID(req.TenantID)
	if err != nil {
		return err
	}

	org, err := c.orgRepository.GetByID(tenant.OrganizationID())
	if err != nil {
		return err
	}

	tenant.Stop()
	pr, err := c.priceRepository.GetByID(tenant.PriceID())
	if err != nil {
		return err
	}

	var balanceFairUsageRefund int64
	if pr.Recurrence() == entities.ProductRecurrenceMonthly {
		activeUntil := time.Unix(tenant.ActiveUntil(), 0)
		activeUntilHours := time.Until(activeUntil).Hours()
		balanceFairUsageRefund = int64(int64(activeUntilHours) / 24 * pr.Price())
	} else {
		activeUntil := time.Unix(tenant.ActiveUntil(), 0)
		activeUntilHours := time.Until(activeUntil).Hours()
		balanceFairUsageRefund = int64(int64(activeUntilHours) / 24 / 365 * pr.Price())
	}

	org.SetBalance(org.Balance() + balanceFairUsageRefund)

	if err := c.orgRepository.Update(org); err != nil {
		return err
	}

	if err := c.tenantRepository.Update(tenant); err != nil {
		return err
	}

	return nil
}
