package services

import "github.com/midtrans/midtrans-go/coreapi"

type Midtrans interface {
	CreateTransaction(orderID string, grossAmount int64) (interface{}, error)
	CheckTransactionStatus(orderID string) (*coreapi.TransactionStatusResponse, error)
}
