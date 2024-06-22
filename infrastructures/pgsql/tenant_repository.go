package pgsql

import (
	"saas-billing/domain/entities"
	"saas-billing/domain/repositories"
	"saas-billing/errors"

	"gorm.io/gorm"
)

type TenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) repositories.TenantRepository {
	return &TenantRepository{db}
}

// GetById finds a tenant by its ID and throw an error if not found
func (tr *TenantRepository) GetById(id string) (*entities.Tenant, error) {
	var tenant entities.Tenant
	if err := tr.db.Where("id = ?", id).First(&tenant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrTenantNotFound
		}
		return nil, err
	}

	if tenant.ID() == "" {
		return nil, errors.ErrTenantNotFound
	}
	return &tenant, nil
}

// Create creates a new tenant
func (tr *TenantRepository) Create(tenant *entities.Tenant) error {
	return tr.db.Create(tenant).Error
}

// Update updates a tenant
func (tr *TenantRepository) Update(tenant *entities.Tenant) error {
	return tr.db.Save(tenant).Error
}
