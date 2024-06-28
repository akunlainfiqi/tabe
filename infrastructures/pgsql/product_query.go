package pgsql

import (
	"saas-billing/app/queries"

	"gorm.io/gorm"
)

type productQuery struct {
	db *gorm.DB
}

func NewProductQuery(db *gorm.DB) queries.ProductQuery {
	return &productQuery{db}
}

type ProductDTO struct {
	PriceID    string
	Price      int64
	Reccurence string
	ProductId  string
	AppID      string
	Name       string
	TierName   string
	TierIndex  int
}

// FindByID finds a product by its ID
func (pq *productQuery) FindByID(id string) (queries.Product, error) {
	var temp []ProductDTO

	if err := pq.db.Raw(`
		SELECT
			pr.id as price_id,
			pr.price,
			pr.reccurence,
			p.id as product_id,
			p.app_id,
			a.name,
			p.tier_name,
			p.tier_index
		FROM
			products p
		INNER JOIN
			prices pr
		ON
			p.id = pr.product_id
		JOIN
			apps a
		ON
			p.app_id = a.id
		WHERE
			p.id = ?
	`, id).Scan(&temp).Error; err != nil {
		return queries.Product{}, err
	}

	product := queries.Product{}
	var price []queries.Price
	for _, t := range temp {
		price = append(price, queries.Price{
			ID:         t.PriceID,
			ProductId:  t.ProductId,
			Price:      t.Price,
			Reccurence: t.Reccurence,
		})

		product = queries.Product{
			ID:        t.PriceID,
			AppId:     t.AppID,
			Name:      t.Name,
			TierName:  t.TierName,
			TierIndex: t.TierIndex,
			Price:     price,
		}
	}

	return product, nil
}

// FindAll finds all products
func (pq *productQuery) FindAll() ([]queries.Product, error) {
	var temp []ProductDTO

	if err := pq.db.Raw(`
		SELECT
			pr.id as price_id,
			pr.price,
			pr.reccurence,
			p.id as product_id,
			p.app_id,
			a.name,
			p.tier_name,
			p.tier_index
		FROM
			products p
		INNER JOIN
			prices pr
		ON
			p.id = pr.product_id
		JOIN
			apps a
		ON
			p.app_id = a.id
	`).Scan(&temp).Error; err != nil {
		return nil, err
	}

	products := make([]queries.Product, 0)
	productExists := false

	for _, t := range temp {
		price := queries.Price{
			ID:         t.PriceID,
			ProductId:  t.ProductId,
			Price:      t.Price,
			Reccurence: t.Reccurence,
		}

		productExists = false
		for i, p := range products {
			if p.ID == t.ProductId {
				products[i].Price = append(products[i].Price, price)
				productExists = true
				break
			}
		}

		if !productExists {
			// Product does not exist, create a new one
			product := queries.Product{
				ID:        t.ProductId,
				AppId:     t.AppID,
				Name:      t.Name,
				TierName:  t.TierName,
				TierIndex: t.TierIndex,
				Price:     []queries.Price{price},
			}

			products = append(products, product)
		}
	}

	return products, nil
}

// FindByAppID finds all products by app ID
func (pq *productQuery) FindByAppID(appID string) ([]queries.Product, error) {
	var temp []ProductDTO

	if err := pq.db.Raw(`
		SELECT
			pr.id as price_id,
			pr.price,
			pr.reccurence,
			p.id as product_id,
			p.app_id,
			a.name,
			p.tier_name,
			p.tier_index
		FROM
			products p
		INNER JOIN
			prices pr
		ON
			p.id = pr.product_id
		JOIN
			apps a
		ON
			p.app_id = a.id
		WHERE
			p.app_id = ?
	`, appID).Scan(&temp).Error; err != nil {
		return nil, err
	}

	products := make([]queries.Product, 0)
	productExists := false

	for _, t := range temp {
		price := queries.Price{
			ID:         t.PriceID,
			ProductId:  t.ProductId,
			Price:      t.Price,
			Reccurence: t.Reccurence,
		}

		productExists = false
		for i, p := range products {
			if p.ID == t.ProductId {
				products[i].Price = append(products[i].Price, price)
				productExists = true
				break
			}
		}

		if !productExists {
			// Product does not exist, create a new one
			product := queries.Product{
				ID:        t.ProductId,
				AppId:     t.AppID,
				Name:      t.Name,
				TierName:  t.TierName,
				TierIndex: t.TierIndex,
				Price:     []queries.Price{price},
			}

			products = append(products, product)
		}
	}

	return products, nil
}
