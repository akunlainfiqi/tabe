package pgsql

import (
	"saas-billing/domain/entities"
	"saas-billing/domain/repositories"
	"saas-billing/errors"

	"gorm.io/gorm"
)

type OrganizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) repositories.OrganizationRepository {
	return &OrganizationRepository{db}
}

// FindByID finds an organization by its ID
func (or *OrganizationRepository) GetByID(id string) (*entities.Organization, error) {
	var dto struct {
		ID             string
		Name           string
		Identifier     string
		Balance        int64
		ContactName    string
		ContactEmail   string
		ContactPhone   string
		ContactAddress string
	}

	if err := or.db.Raw(`
		SELECT
			id,
			name,
			identifier,
			balance,
			contact_name,
			contact_email,
			contact_phone,
			contact_address
		FROM
			organizations
		WHERE
			id = @id
	`, map[string]interface{}{
		"id": id,
	}).Scan(&dto); err != nil {
		return nil, err.Error
	}

	if dto.ID == "" {
		return nil, errors.ErrOrganizationNotFound
	}

	organization := entities.BuildOrganization(dto.ID, dto.Name, dto.Identifier, dto.Balance, dto.ContactName, dto.ContactEmail, dto.ContactPhone, dto.ContactAddress)

	return organization, nil
}

// Create creates a new organization
func (or *OrganizationRepository) Create(organization *entities.Organization) error {
	err := or.db.Exec(`
		INSERT INTO organizations (id, name, identifier, balance, contact_name, contact_email, contact_phone, contact_address)
		VALUES (@id, @name, @identifier, @balance, @contact_name, @contact_email, @contact_phone, @contact_address)
	`, map[string]interface{}{
		"id":              organization.ID(),
		"name":            organization.Name(),
		"identifier":      organization.Identifier(),
		"balance":         organization.Balance(),
		"contact_name":    organization.ContactName(),
		"contact_email":   organization.ContactEmail(),
		"contact_phone":   organization.ContactPhone(),
		"contact_address": organization.ContactAddress(),
	}).Error
	return err
}

// Update updates an organization
func (or *OrganizationRepository) Update(organization *entities.Organization) error {
	err := or.db.Exec(`
		UPDATE organizations
		SET name = @name, identifier = @identifier, balance = @balance, contact_name = @contact_name, contact_email = @contact_email, contact_phone = @contact_phone, contact_address = @contact_address
		WHERE id = @id
	`, map[string]interface{}{
		"id":              organization.ID(),
		"name":            organization.Name(),
		"identifier":      organization.Identifier(),
		"balance":         organization.Balance(),
		"contact_name":    organization.ContactName(),
		"contact_email":   organization.ContactEmail(),
		"contact_phone":   organization.ContactPhone(),
		"contact_address": organization.ContactAddress(),
	}).Error
	return err
}
