package pgsql

import (
	"saas-billing/domain/entities"
	"saas-billing/domain/repositories"

	"gorm.io/gorm"
)

type OrganizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) repositories.OrganizationRepository {
	return &OrganizationRepository{db}
}

// FindByID finds an organization by its ID
func (or *OrganizationRepository) FindByID(id string) (*entities.Organization, error) {
	var organization entities.Organization
	if err := or.db.Where("id = ?", id).First(&organization).Error; err != nil {
		return nil, err
	}
	return &organization, nil
}

// Create creates a new organization
func (or *OrganizationRepository) Create(organization *entities.Organization) error {
	return or.db.Create(organization).Error
}

// Update updates an organization
func (or *OrganizationRepository) Update(organization *entities.Organization) error {
	return or.db.Save(organization).Error
}
