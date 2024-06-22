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
	var iamOrganization entities.IamOrganization
	if err := ior.db.Where("id = ?", id).First(&iamOrganization).Error; err != nil {
		return nil, err
	}
	if iamOrganization.ID == "" {
		return nil, errors.ErrOrganizationNotFound
	}
	return &iamOrganization, nil
}
