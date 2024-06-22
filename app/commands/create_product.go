package commands

import (
	"saas-billing/domain/entities"
	"saas-billing/domain/repositories"

	"github.com/google/uuid"
)

type CreateProductRequest struct {
	name  string
	tiers []CreateProductTierRequest
}

type CreateProductTierRequest struct {
	Name  string                      `json:"name"`
	Order int                         `json:"order"`
	Price []CreateProductPriceRequest `json:"price"`
}

type CreateProductPriceRequest struct {
	Price      float64 `json:"price"`
	Reccurence string  `json:"reccurence"`
}

func NewCreateProductRequest(
	name string,
	tiers []CreateProductTierRequest,
) *CreateProductRequest {
	return &CreateProductRequest{
		name:  name,
		tiers: tiers,
	}
}

type CreateProductCommand struct {
	productRepository repositories.ProductRepository
	appsRepository    repositories.AppsRepository
	priceRepository   repositories.PriceRepository
}

func NewCreateProductCommand(
	productRepository repositories.ProductRepository,
	appsRepository repositories.AppsRepository,
	priceRepository repositories.PriceRepository,
) *CreateProductCommand {
	return &CreateProductCommand{
		productRepository: productRepository,
		appsRepository:    appsRepository,
		priceRepository:   priceRepository,
	}
}

func (c *CreateProductCommand) Execute(req *CreateProductRequest) error {
	app := entities.NewApps(uuid.New().String(), req.name)

	if err := c.appsRepository.Create(app); err != nil {
		return err
	}

	for _, tier := range req.tiers {
		product := entities.NewProduct(uuid.New().String(), *app, tier.Name, tier.Order)

		if err := c.productRepository.Create(product); err != nil {
			return err
		}

		for _, price := range tier.Price {
			productPrice := entities.NewPrice(uuid.New().String(), *product, price.Price, price.Reccurence)

			if err := c.priceRepository.Create(productPrice); err != nil {
				return err
			}
		}
	}

	return nil
}
