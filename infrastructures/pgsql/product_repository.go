package pgsql

import (
	"saas-billing/domain/entities"
	"saas-billing/domain/repositories"
	"saas-billing/errors"

	"gorm.io/gorm"
)

// ProductRepository is a repository for managing Product entities
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new ProductRepository
func NewProductRepository(db *gorm.DB) repositories.ProductRepository {
	return &ProductRepository{db}
}

// FindByID finds a product by its ID
func (pr *ProductRepository) FindByID(id string) (*entities.Product, error) {
	var dto struct {
		ProductId string
		AppID     string
		Name      string
		TierName  string
		TierIndex int
	}

	if err := pr.db.Raw(`
		SELECT
			p.id,
			p.app_id,
			a.name,
			p.tier_name,
			p.tier_index
		FROM
			products p
		JOIN
			apps a
		ON
			p.app_id = a.id
		WHERE
			p.id = @id
	`, map[string]interface{}{
		"id": id,
	}).Scan(&dto); err != nil {
		return nil, err.Error
	}

	if dto.ProductId == "" {
		return nil, errors.ErrProductNotFound
	}

	app := entities.NewApps(dto.AppID, dto.Name)
	product := entities.NewProduct(dto.ProductId, *app, dto.TierName, dto.TierIndex)

	return product, nil
}

// Create creates a new product
func (pr *ProductRepository) Create(product *entities.Product) error {
	if err := pr.db.Exec(`
		INSERT INTO products (id, app_id, tier_name, tier_index)
		VALUES (@id, @app_id, @tier_name, @tier_index)
	`, map[string]interface{}{
		"id":         product.ID(),
		"app_id":     product.App().ID(),
		"tier_name":  product.TierName(),
		"tier_index": product.TierIndex(),
	}).Error; err != nil {
		return err
	}

	return nil
}
