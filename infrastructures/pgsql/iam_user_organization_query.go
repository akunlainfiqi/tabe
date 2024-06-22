package pgsql

import (
	"saas-billing/app/queries"

	"gorm.io/gorm"
)

type IamUserOrganizationQuery struct {
	db *gorm.DB
}

func NewIamUserOrganizationQuery(db *gorm.DB) queries.IamUserOrganizationQuery {
	return &IamUserOrganizationQuery{db}
}

func (q *IamUserOrganizationQuery) IsOwner(organizationID, userID string) bool {
	var count int64
	q.db.Raw(`
		SELECT
			COUNT(*)
		FROM
			user_organization
		WHERE
			organization_id = ?
		AND
			user_id = ?
		AND
			level = 'owner'
	`, organizationID, userID).Scan(&count)

	return count > 0
}
