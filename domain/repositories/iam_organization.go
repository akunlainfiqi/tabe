package repositories

import (
	"saas-billing/domain/entities"
)

type IamOrganizationRepository interface {
	GetByID(id string) (*entities.IamOrganization, error)
	IsOwner(organizationID, userID string) bool
}
