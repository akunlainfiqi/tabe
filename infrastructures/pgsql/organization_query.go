package pgsql

import (
	"saas-billing/app/queries"
	"saas-billing/errors"

	"gorm.io/gorm"
)

type organizationQuery struct {
	db *gorm.DB
}

func NewOrganizationQuery(db *gorm.DB) queries.OrganizationQuery {
	return &organizationQuery{db}
}

func (oq *organizationQuery) GetByID(id string) (queries.Organization, error) {
	var temp []queries.Organization
	if err := oq.db.Raw(`
		SELECT
			id,
			name,
			identifier,
			balance
		FROM
			organizations
		WHERE
			id = ?
	`, id).Scan(&temp).Error; err != nil {
		return queries.Organization{}, err
	}

	if len(temp) == 0 {
		return queries.Organization{}, errors.ErrOrganizationNotFound
	}

	return temp[0], nil
}
