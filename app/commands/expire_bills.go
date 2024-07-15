package commands

import (
	"saas-billing/domain/repositories"
)

type ExpireBillsCommand struct {
	billsRepository repositories.BillsRepository
}

func NewExpireBillsCommand(billsRepository repositories.BillsRepository) *ExpireBillsCommand {
	return &ExpireBillsCommand{
		billsRepository: billsRepository,
	}
}

func (c *ExpireBillsCommand) Execute() error {
	bills, err := c.billsRepository.GetUnpaidBillsAfterDueDate()
	if err != nil {
		return err
	}

	for _, bill := range bills {
		bill.Expire()

		if err := c.billsRepository.Update(bill); err != nil {
			return err
		}
	}

	return nil
}
