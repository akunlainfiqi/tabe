package pgsql

import (
	"saas-billing/app/queries"

	"gorm.io/gorm"
)

type AppQuery struct {
	db *gorm.DB
}

func NewAppQuery(db *gorm.DB) queries.AppQueries {
	return &AppQuery{db}
}

func (q *AppQuery) GetAll() ([]queries.App, error) {
	var apps []queries.App

	if err := q.db.Raw(`
		SELECT
			id,
			name,
			image_url
		FROM
			apps
	`).Scan(&apps).Error; err != nil {
		return nil, err
	}

	return apps, nil
}
