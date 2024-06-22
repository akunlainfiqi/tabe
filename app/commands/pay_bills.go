package commands

import (
	"saas-billing/domain/entities"
	"saas-billing/domain/repositories"

	"github.com/google/uuid"
)

// TODO: implement payment gateway services
type PayBillsRequest struct {
	billId               string
	transactionType      string
	transactionTimestamp int64
}

func NewPayBillsRequest(billId, transactionType string, transactionTimestamp int64) *PayBillsRequest {
	return &PayBillsRequest{
		billId:               billId,
		transactionType:      transactionType,
		transactionTimestamp: transactionTimestamp,
	}
}

type PayBillsCommand struct {
	billsRepository       repositories.BillsRepository
	tenantRepository      repositories.TenantRepository
	transactionRepository repositories.TransactionRepository
}

func NewPayBillsCommand(
	billsRepository repositories.BillsRepository,
	tenantRepository repositories.TenantRepository,
	transactionRepository repositories.TransactionRepository,
) *PayBillsCommand {
	return &PayBillsCommand{
		billsRepository:       billsRepository,
		tenantRepository:      tenantRepository,
		transactionRepository: transactionRepository,
	}
}

func (c *PayBillsCommand) Execute(req *PayBillsRequest) error {
	bills, err := c.billsRepository.GetByID(req.billId)
	if err != nil {
		return err
	}

	tenant, err := c.tenantRepository.GetById(bills.TenantID())
	if err != nil {
		return err
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	transaction, err := entities.NewTransaction(
		id.String(),
		bills.ID(),
		tenant.OrganizationID(),
		bills.Amount(),
		req.transactionType,
		req.transactionTimestamp,
	)
	if err != nil {
		return err
	}

	if err := c.transactionRepository.Create(transaction); err != nil {
		return err
	}

	bills.SetStatus(entities.BillStatusPaid)

	if err := c.billsRepository.Update(bills); err != nil {
		return err
	}
	return nil
}
