package midtransapi

import (
	"saas-billing/app/services"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type midtransService struct {
	CoreClient *coreapi.Client
	SnapClient *snap.Client
}

func NewMidtrans(core_client *coreapi.Client, snap *snap.Client) services.Midtrans {
	return &midtransService{
		CoreClient: core_client,
		SnapClient: snap,
	}
}

func (m *midtransService) CreateTransaction(orderID string, grossAmount int64) (interface{}, error) {
	res, err := m.SnapClient.CreateTransaction(&snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: grossAmount,
		},
		Expiry: &snap.ExpiryDetails{
			Unit:     "hour",
			Duration: 1,
		},
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *midtransService) CheckTransactionStatus(orderID string) (*coreapi.TransactionStatusResponse, error) {
	res, err := m.CoreClient.CheckTransaction(orderID)
	if err != nil {
		return nil, err
	}
	return res, nil
}
