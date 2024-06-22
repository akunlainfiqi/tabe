package repositories

import "saas-billing/domain/entities"

type AppsRepository interface {
	FindByID(id string) (*entities.Apps, error)
	Create(apps *entities.Apps) error
}
