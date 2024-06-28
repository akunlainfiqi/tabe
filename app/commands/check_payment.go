package commands

import (
	"fmt"
	"saas-billing/app/services"
	"saas-billing/domain/entities"
	"saas-billing/domain/repositories"
	"time"

	"github.com/google/uuid"
)

type CheckPaymentCommand struct {
	transactionRepository repositories.TransactionRepository
	billRepository        repositories.BillsRepository
	orgRepository         repositories.OrganizationRepository

	midtransService services.Midtrans
}

func NewCheckPayment(
	transactionRepository repositories.TransactionRepository,
	billRepository repositories.BillsRepository,
	orgRepository repositories.OrganizationRepository,
	midtransService services.Midtrans,
) *CheckPaymentCommand {
	return &CheckPaymentCommand{
		transactionRepository: transactionRepository,
		billRepository:        billRepository,
		orgRepository:         orgRepository,
		midtransService:       midtransService,
	}
}

func parseTransactionTime(transactionTime string) (int64, error) {
	const layout = "2006-01-02 15:04:05"

	location, err := time.LoadLocation("Asia/Bangkok") // UTC+7 biasanya menggunakan Asia/Bangkok
	if err != nil {
		return 0, fmt.Errorf("failed to load location: %v", err)
	}

	// Parse string waktu menjadi objek waktu dengan lokasi UTC+7
	parsedTime, err := time.ParseInLocation(layout, transactionTime, location)
	if err != nil {
		return 0, err
	}
	timeInUTC := parsedTime.UTC()

	return timeInUTC.Unix(), nil
}

func (c *CheckPaymentCommand) checkBillById(bill *entities.Bills) error {
	res, err := c.midtransService.CheckTransactionStatus(bill.ID())
	if err != nil {
		bill.Cancel()
		if err := c.billRepository.Update(bill); err != nil {
			return err
		}

		return nil
	}

	if res.TransactionStatus == "settlement" {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}

		trtime, err := parseTransactionTime(res.TransactionTime)
		if err != nil {
			return err
		}

		bill.Settle()
		if err := c.billRepository.Update(bill); err != nil {
			return err
		}

		tr, err := entities.NewTransaction(
			id.String(),
			bill.ID(),
			bill.OrganizationID(),
			bill.Amount(),
			entities.TransactionTypePayment,
			trtime,
		)

		if err != nil {
			return err
		}

		if err := c.transactionRepository.Create(tr); err != nil {
			return err
		}

	} else if res.TransactionStatus == "expire" {
		if bill.BalanceUsed() > 0 {
			org, err := c.orgRepository.GetByID(bill.OrganizationID())
			if err != nil {
				return err
			}

			org.SetBalance(org.Balance() + bill.BalanceUsed())
			if err := c.orgRepository.Update(org); err != nil {
				return err
			}
		}
		bill.Expire()
		if err := c.billRepository.Update(bill); err != nil {
			return err
		}
	} else if res.TransactionStatus == "cancel" {
		if bill.BalanceUsed() > 0 {
			org, err := c.orgRepository.GetByID(bill.OrganizationID())
			if err != nil {
				return err
			}

			org.SetBalance(org.Balance() + bill.BalanceUsed())
			if err := c.orgRepository.Update(org); err != nil {
				return err
			}
		}

		bill.Cancel()
		if err := c.billRepository.Update(bill); err != nil {
			return err
		}

	}

	return nil
}

func (c *CheckPaymentCommand) CheckBillById(billId string) error {
	bill, err := c.billRepository.GetByID(billId)
	if err != nil {
		return err
	}

	return c.checkBillById(bill)
}

func (c *CheckPaymentCommand) Execute() error {
	bills, err := c.billRepository.GetUnpaidBillsAfterDueDate()
	if err != nil {
		return err
	}

	for _, bill := range bills {
		if err := c.checkBillById(bill); err != nil {
			return err
		}
	}

	return nil
}