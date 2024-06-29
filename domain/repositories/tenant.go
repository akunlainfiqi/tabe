package repositories

import "saas-billing/domain/entities"

type TenantRepository interface {
	GetByID(id string) (*entities.Tenant, error)
	Create(tenant *entities.Tenant) error
	Update(tenant *entities.Tenant) error
}
