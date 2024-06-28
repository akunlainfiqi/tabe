package repositories

import "saas-billing/domain/entities"

type OrganizationRepository interface {
	GetByID(id string) (*entities.Organization, error)
	FindByID(id string) (*entities.Organization, error)
	Create(organization *entities.Organization) error
	Update(organization *entities.Organization) error
}
