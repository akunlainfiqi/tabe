package repositories

import "saas-billing/domain/entities"

type TenantRepository interface {
	GetById(id string) (*entities.Tenant, error)
	Create(tenant *entities.Tenant) error
	Update(tenant *entities.Tenant) error
}
