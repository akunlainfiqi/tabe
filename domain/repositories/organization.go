package repositories

import "saas-billing/domain/entities"

type OrganizationRepository interface {
	FindByID(id string) (*entities.Organization, error)
	Create(organization *entities.Organization) error
	Update(organization *entities.Organization) error
}