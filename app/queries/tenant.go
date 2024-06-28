package queries

type Tenant struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ProductID   string `json:"product_id"`
	OrgID       string `json:"org_id"`
	ActiveUntil string `json:"active_until"`
	CreatedAt   string `json:"created_at"`
}

type TenantQuery interface {
	FindByID(id string) (Tenant, error)
	FindByOrgID(orgID string) ([]Tenant, error)
}
