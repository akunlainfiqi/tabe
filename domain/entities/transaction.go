package entities

type Transaction struct {
	id                   string
	billsId              string
	organizationId       string
	amount               int64
	transactionType      string
	transactionTimestamp int64
}

var (
	TransactionTypeBalance = "balance"
)

func BuildTransaction(
	id,
	billsId,
	organizationId string,
	amount int64,
	transactionType string,
	transactionTimestamp int64,
) *Transaction {
	return &Transaction{
		id:                   id,
		billsId:              billsId,
		organizationId:       organizationId,
		amount:               amount,
		transactionType:      transactionType,
		transactionTimestamp: transactionTimestamp,
	}
}

func NewTransaction(
	id,
	billsId,
	organizationId string,
	amount int64,
	transactionType string,
	transactionTimestamp int64,
) (*Transaction, error) {
	return &Transaction{
		id:                   id,
		billsId:              billsId,
		organizationId:       organizationId,
		amount:               amount,
		transactionType:      transactionType,
		transactionTimestamp: transactionTimestamp,
	}, nil
}

func (t *Transaction) ID() string {
	return t.id
}

func (t *Transaction) BillsID() string {
	return t.billsId
}

func (t *Transaction) OrganizationID() string {
	return t.organizationId
}

func (t *Transaction) Amount() int64 {
	return t.amount
}

func (t *Transaction) TransactionType() string {
	return t.transactionType
}

func (t *Transaction) TransactionTimestamp() int64 {
	return t.transactionTimestamp
}
