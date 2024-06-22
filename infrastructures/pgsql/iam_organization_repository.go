package pgsql

import (
	"saas-billing/domain/entities"
	"saas-billing/domain/repositories"
	"saas-billing/errors"

	"gorm.io/gorm"
)

type IamOrganizationRepository struct {
	db *gorm.DB
}

func NewIamOrganizationRepository(db *gorm.DB) repositories.IamOrganizationRepository {
	return &IamOrganizationRepository{db}
}

// GetByID finds an IAM organization by its ID
func (ior *IamOrganizationRepository) GetByID(id string) (*entities.IamOrganization, error) {
	var temp struct {
		ID         string
		Name       string
		Identifier string
	}
	ior.db.Raw(`
	SELECT *
	FROM organization
	WHERE id = ?
		`, id).Scan(&temp)
	if temp.ID == "" {
		return nil, errors.ErrOrganizationNotFound
	}
	iamOrganization := entities.IamOrganization{
		ID:         temp.ID,
		Name:       temp.Name,
		Identifier: temp.Identifier,
	}
	return &iamOrganization, nil
}

// IsOwner checks if a user is the owner of an organization
func (ior *IamOrganizationRepository) IsOwner(organizationID, userID string) bool {
	var count int64
	ior.db.Raw(`
		SELECT COUNT(*)
		FROM user_organization
		WHERE organization_id = ? AND user_id = ? AND level = 'owner'
		`, organizationID, userID).Count(&count)
	return count > 0
}
