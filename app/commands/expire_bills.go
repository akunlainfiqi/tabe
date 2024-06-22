package commands

import (
	"saas-billing/domain/repositories"
)

type ExpireBillsRequest struct {
	billId string
}

func NewExpireBillsRequest(billId string) *ExpireBillsRequest {
	return &ExpireBillsRequest{
		billId: billId,
	}
}

type ExpireBillsCommand struct {
	billsRepository repositories.BillsRepository
}

func NewExpireBillsCommand(billsRepository repositories.BillsRepository) *ExpireBillsCommand {
	return &ExpireBillsCommand{
		billsRepository: billsRepository,
	}
}

func (c *ExpireBillsCommand) Execute(req *ExpireBillsRequest) error {
	bills, err := c.billsRepository.GetByID(req.billId)
	if err != nil {
		return err
	}

	bills.Expire()

	return c.billsRepository.Update(bills)
}
