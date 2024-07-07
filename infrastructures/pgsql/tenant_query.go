package pgsql

import (
	"saas-billing/app/queries"

	"gorm.io/gorm"
)

type tenantQuery struct {
	db *gorm.DB
}

func NewTenantQuery(db *gorm.DB) queries.TenantQuery {
	return &tenantQuery{db}
}

type TenantDTO struct {
	ID             string
	Name           string
	AppID          string
	AppName        string
	OrganizationId string
	ActiveUntil    string
	CreatedAt      string
}

// FindByID finds a tenant by its ID
func (tq *tenantQuery) FindByID(id string) (queries.Tenant, error) {
	var temp []TenantDTO

	if err := tq.db.Raw(`
		SELECT
			t.id,
			t.name,
			t.app_id,
			t.organization_id,
			t.active_until,
			t.created_at
		FROM
			tenants t
		WHERE
			t.id = ?
	`, id).Scan(&temp).Error; err != nil {
		return queries.Tenant{}, err
	}

	tenant := queries.Tenant{}
	for _, t := range temp {
		tenant.ID = t.ID
		tenant.Name = t.Name
		tenant.AppID = t.AppID
		tenant.OrgID = t.OrganizationId
		tenant.ActiveUntil = t.ActiveUntil
		tenant.CreatedAt = t.CreatedAt
	}

	return tenant, nil
}

// FindByOrgID finds a tenant by its OrgID
func (tq *tenantQuery) FindByOrgID(orgID string) ([]queries.Tenant, error) {
	var temp []TenantDTO

	if err := tq.db.Raw(`
		SELECT
			t.id,
			t.name,
			t.app_id,
			a.name as app_name,
			t.organization_id,
			t.active_until,
			t.created_at
		FROM
			tenants t
		JOIN
			apps a
		ON
			t.app_id = a.id
		WHERE
			t.organization_id = ?
	`, orgID).Scan(&temp).Error; err != nil {
		return []queries.Tenant{}, err
	}

	tenants := []queries.Tenant{}
	for _, t := range temp {
		tenant := queries.Tenant{}
		tenant.ID = t.ID
		tenant.Name = t.Name
		tenant.AppID = t.AppID
		tenant.AppName = t.AppName
		tenant.OrgID = t.OrganizationId
		tenant.ActiveUntil = t.ActiveUntil
		tenant.CreatedAt = t.CreatedAt
		tenants = append(tenants, tenant)
	}

	return tenants, nil
}
