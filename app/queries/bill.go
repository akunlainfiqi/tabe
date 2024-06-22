package queries

type Bill struct {
	Id             string `json:"id"`
	OrganizationId string `json:"organization_id"`
	TenantId       string `json:"tenant_id"`
	TenantName     string `json:"tenant_name"`
	Status         string `json:"status"`
	DueDate        int64  `json:"due_date"`
	Amount         int64  `json:"amount"`
	BalanceUsed    int64  `json:"balance_used"`
	CreatedAt      int64  `json:"created_at"`
}

type BillQuery interface {
	GetByID(id string) (Bill, error)
	GetByOrganizationID(organizationID string) ([]Bill, error)
}
