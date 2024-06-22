package repositories

import "saas-billing/domain/entities"

type RoleRepository interface {
	FindByID(id string) (*entities.Role, error)
	FindManyByUserID(id string) ([]entities.Role, error)
	Create(role *entities.Role) error
}

