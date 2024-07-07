package queries

type Transaction struct {
	ID                   string `json:"id"`
	BillsID              string `json:"bills_id"`
	OrganizationID       string `json:"organization_id"`
	Amount               int64  `json:"amount"`
	TransactionType      string `json:"transaction_type"`
	TransactionTimestamp int64  `json:"transaction_timestamp"`
}

type TransactionQuery interface {
	FindByOrgID(id string) ([]Transaction, error)
	FindByBillsID(billsID string) (Transaction, error)
}
