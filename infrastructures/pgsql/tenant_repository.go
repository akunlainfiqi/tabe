package pgsql

import (
	"saas-billing/domain/entities"
	"saas-billing/domain/repositories"
	"saas-billing/errors"
	"time"

	"gorm.io/gorm"
)

type TenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) repositories.TenantRepository {
	return &TenantRepository{db}
}

// GetById finds a tenant by its ID and throw an error if not found
func (tr *TenantRepository) GetByID(id string) (*entities.Tenant, error) {
	var dto struct {
		ID             string
		Name           string
		AppID          string
		OrganizationId string
		ActiveUntil    int64
		PriceId        string
	}

	tr.db.Raw(`
		SELECT
			t.id,
			t.name,
			t.app_id,
			t.organization_id,
			t.active_until,
			t.price_id
		FROM tenants t
		WHERE t.id = ?
	`, id).Scan(&dto)

	tenant := entities.BuildTenant(dto.ID, dto.Name, dto.AppID, dto.OrganizationId, dto.PriceId, dto.ActiveUntil)

	if tenant.ID() == "" {
		return nil, errors.ErrTenantNotFound
	}
	return tenant, nil
}

// Create creates a new tenant
func (tr *TenantRepository) Create(tenant *entities.Tenant) error {
	err := tr.db.Exec(`
		INSERT INTO 
			tenants (
				id, 
				name, 
				app_id, 
				organization_id, 
				price_id, 
				active_until, 
				created_at, 
				updated_at
				)
		VALUES (
			@id, 
			@name, 
			@app_id, 
			@organization_id, 
			@price_id, 
			@active_until, 
			@date, 
			@date
			)
	`, map[string]interface{}{
		"id":              tenant.ID(),
		"name":            tenant.Name(),
		"app_id":          tenant.AppID(),
		"organization_id": tenant.OrganizationID(),
		"price_id":        tenant.PriceID(),
		"active_until":    tenant.ActiveUntil(),
		"date":            time.Now().Unix(),
	}).Error
	return err
}

// Update updates a tenant
func (tr *TenantRepository) Update(tenant *entities.Tenant) error {
	err := tr.db.Exec(`
		UPDATE 
			tenants
		SET 
			name = @name, 
			app_id = @app_id, 
			organization_id = @organization_id, 
			price_id = @price_id, 
			active_until = @active_until, 
			updated_at = @date
		WHERE 
			id = @id
	`, map[string]interface{}{
		"id":              tenant.ID(),
		"name":            tenant.Name(),
		"app_id":          tenant.AppID(),
		"organization_id": tenant.OrganizationID(),
		"price_id":        tenant.PriceID(),
		"active_until":    tenant.ActiveUntil(),
		"date":            time.Now().Unix(),
	}).Error
	return err
}
