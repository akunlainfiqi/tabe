package commands

import (
	"saas-billing/domain/entities"
	"saas-billing/domain/repositories"
	"time"

	"github.com/google/uuid"
)

// CreateBillsRequest is a request for creating bills
type CreateBillsRequest struct {
	tenantId string
	billType string
}

func NewCreateBillsRequest(
	tenantId string,
	billType string,
) *CreateBillsRequest {
	return &CreateBillsRequest{
		tenantId: tenantId,
		billType: billType,
	}
}

// CreateBillsCommand is a command for creating bills
type CreateBillsCommand struct {
	billsRepository        repositories.BillsRepository
	tenantRepository       repositories.TenantRepository
	organizationRepository repositories.OrganizationRepository
	transactionRepository  repositories.TransactionRepository
	priceRepository        repositories.PriceRepository
}

// NewCreateBillsCommand creates a new create bills command
func NewCreateBillsCommand(
	billsRepository repositories.BillsRepository,
	tenantRepository repositories.TenantRepository,
	organizationRepository repositories.OrganizationRepository,
	transactionRepository repositories.TransactionRepository,
	priceRepository repositories.PriceRepository,
) *CreateBillsCommand {
	return &CreateBillsCommand{
		billsRepository:        billsRepository,
		tenantRepository:       tenantRepository,
		organizationRepository: organizationRepository,
		transactionRepository:  transactionRepository,
		priceRepository:        priceRepository,
	}
}

// Execute executes the create bills command
func (c *CreateBillsCommand) Execute(req *CreateBillsRequest) error {
	tenant, err := c.tenantRepository.GetById(req.tenantId)
	if err != nil {
		return err
	}

	organization, err := c.organizationRepository.GetByID(tenant.OrganizationID())
	if err != nil {
		return err
	}

	billId, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	price, err := c.priceRepository.GetByID(tenant.PriceID())
	if err != nil {
		return err
	}
	balanceUsed := int64(0)

	if organization.Balance() > 0 {
		if organization.Balance() < price.Price() {
			balanceUsed = organization.Balance()
		} else {
			balanceUsed = price.Price()
		}
	}

	bill, err := entities.NewBills(
		billId.String(),
		organization.ID(),
		tenant.ID(),
		price.Price(),
		balanceUsed,
		time.Now().AddDate(0, 0, 10).Unix(),
		req.billType,
	)
	if err != nil {
		return err
	}

	if err := c.billsRepository.Create(bill); err != nil {
		return err
	}

	organization.SetBalance(organization.Balance() - balanceUsed)

	if err := c.organizationRepository.Update(organization); err != nil {
		return err
	}

	if bill.Status() == entities.BillStatusPaid {
		transactionId, err := uuid.NewRandom()
		if err != nil {
			return err
		}

		transaction, err := entities.NewTransaction(
			transactionId.String(),
			bill.ID(),
			organization.ID(),
			balanceUsed,
			entities.TransactionTypeBalance,
			time.Now().Unix(),
		)
		if err != nil {
			return err
		}

		if err := c.transactionRepository.Create(transaction); err != nil {
			return err
		}

		var activeUntil int64
		if tenant.ActiveUntil() == 0 || tenant.ActiveUntil() < time.Now().Unix() {
			if price.Recurrence() == entities.ProductRecurrenceMonthly {
				activeUntil = time.Now().AddDate(0, 1, 0).Unix()
			}
			if price.Recurrence() == entities.ProductRecurrenceYearly {
				activeUntil = time.Now().AddDate(1, 0, 0).Unix()
			}
		} else if tenant.ActiveUntil() > time.Now().Unix() {
			if price.Recurrence() == entities.ProductRecurrenceMonthly {
				activeUntil = time.Unix(tenant.ActiveUntil(), 0).AddDate(0, 1, 0).Unix()
			}
			if price.Recurrence() == entities.ProductRecurrenceYearly {
				activeUntil = time.Unix(tenant.ActiveUntil(), 0).AddDate(1, 0, 0).Unix()
			}
		}

		tenant.SetActiveUntil(activeUntil)

		if err := c.tenantRepository.Update(tenant); err != nil {
			return err
		}
	}

	return nil
}
