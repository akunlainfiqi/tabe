package pgsql

import (
	"log"
	"saas-billing/domain/entities"
	"saas-billing/domain/repositories"
	"saas-billing/errors"

	"gorm.io/gorm"
)

type PriceRepository struct {
	db *gorm.DB
}

func NewPriceRepository(db *gorm.DB) repositories.PriceRepository {
	return &PriceRepository{db}
}

// FindByID finds a price by its ID
func (pr *PriceRepository) GetByID(id string) (*entities.Price, error) {
	var dto struct {
		ID         string
		ProductId  string
		AppId      string
		AppName    string
		Price      int64
		Reccurence string
		TierName   string
		TierIndex  int
	}
	pr.db.Raw(`
		SELECT 
			p.id,
			p.product_id,
			pr.app_id,
			a.name AS app_name,
			p.price,
			p.reccurence,
			pr.tier_name,
			pr.tier_index
		FROM prices p
		JOIN products pr ON p.product_id = pr.id
		JOIN apps a ON pr.app_id = a.id
		WHERE p.id = ?
		`, id).Scan(&dto)

	log.Print(dto)

	app := entities.NewApps(dto.AppId, dto.AppName)
	product := entities.NewProduct(dto.ProductId, *app, dto.TierName, dto.TierIndex)
	price := entities.NewPrice(dto.ID, *product, dto.Price, dto.Reccurence)

	if price.ID() == "" {
		return nil, errors.ErrPriceNotFound
	}
	return price, nil
}

// Create creates a new price
func (pr *PriceRepository) Create(price *entities.Price) error {
	return pr.db.Create(price).Error
}
