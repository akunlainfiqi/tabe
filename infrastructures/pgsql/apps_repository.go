package pgsql

import (
	"saas-billing/domain/entities"
	"saas-billing/domain/repositories"
	"saas-billing/errors"

	"gorm.io/gorm"
)

type AppsRepository struct {
	db *gorm.DB
}

func NewAppsRepository(db *gorm.DB) repositories.AppsRepository {
	return &AppsRepository{db}
}

// FindByID finds an app by its ID
func (ar *AppsRepository) FindByID(id string) (*entities.Apps, error) {
	var dto struct {
		ID   string
		Name string
	}
	if err := ar.db.Raw(`
		SELECT
			id,
			name
		FROM
			apps
		WHERE
			id = @id
	`, map[string]interface{}{
		"id": id,
	}).Scan(&dto); err != nil {
		return nil, err.Error
	}

	if dto.ID == "" {
		return nil, errors.ErrAppsNotFound
	}

	app := entities.NewApps(dto.ID, dto.Name)

	return app, nil
}

// Create creates a new app
func (ar *AppsRepository) Create(app *entities.Apps) error {
	err := ar.db.Exec(`
		INSERT INTO apps (id, name)
		VALUES (@id, @name)
	`, map[string]interface{}{
		"id":   app.ID(),
		"name": app.Name(),
	}).Error
	if err != nil {
		return err
	}
	return nil
}
