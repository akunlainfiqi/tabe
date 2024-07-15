package commands

import (
	"saas-billing/app/services"
	"saas-billing/domain/entities"
	"saas-billing/domain/repositories"
	"saas-billing/errors"
)

type CreatePaymentURLRequest struct {
	BillID string
}

type CreatePaymentURLCommand struct {
	billsRepository repositories.BillsRepository
	midtransService services.Midtrans
}

func NewCreatePaymentURLCommand(
	billsRepository repositories.BillsRepository,
	midtransService services.Midtrans,
) *CreatePaymentURLCommand {
	return &CreatePaymentURLCommand{
		billsRepository: billsRepository,
		midtransService: midtransService,
	}
}

func (c *CreatePaymentURLCommand) Execute(req *CreatePaymentURLRequest) (interface{}, error) {
	bill, err := c.billsRepository.GetByID(req.BillID)
	if err != nil {
		return nil, err
	}

	if bill.Status() != entities.BillStatusWaitingPayment {
		return nil, errors.ErrInvalidBillStatus
	}

	paymentURL, err := c.midtransService.CreateTransaction(bill.ID(), bill.Total())
	if err != nil {
		return nil, err
	}

	return paymentURL, nil
}
