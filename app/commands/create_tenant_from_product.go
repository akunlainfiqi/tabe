package commands

import (
	"errors"
	"saas-billing/domain/entities"
	"saas-billing/domain/repositories"
)

type CreateTenantFromProductRequest struct {
	priceId        string
	organizationId string
	tenantId       string
	tenantName     string
}

func NewCreateTenantFromProductRequest(
	priceId,
	organizationId,
	tenantId,
	tenantName string,
) (*CreateTenantFromProductRequest, error) {
	return &CreateTenantFromProductRequest{
		priceId:        priceId,
		organizationId: organizationId,
		tenantId:       tenantId,
		tenantName:     tenantName,
	}, nil
}

type CreateTenantFromProductCommand struct {
	tenantRepository       repositories.TenantRepository
	organizationRepository repositories.OrganizationRepository
	priceRepository        repositories.PriceRepository
}

func NewCreateTenantFromProductCommand(
	tenantRepository repositories.TenantRepository,
	organizationRepository repositories.OrganizationRepository,
	priceRepository repositories.PriceRepository,
) *CreateTenantFromProductCommand {
	return &CreateTenantFromProductCommand{
		tenantRepository:       tenantRepository,
		organizationRepository: organizationRepository,
		priceRepository:        priceRepository,
	}
}

func (c *CreateTenantFromProductCommand) Execute(req *CreateTenantFromProductRequest) error {
	organization, err := c.organizationRepository.FindByID(req.organizationId)
	if err != nil {
		return err
	}

	if organization == nil {
		return errors.New("organization not found")
	}

	price, err := c.priceRepository.FindByID(req.priceId)
	if err != nil {
		return err
	}

	if price == nil {
		return errors.New("price not found")
	}

	product := price.Product()

	tenant := entities.NewTenant(req.tenantId, req.tenantName, product.ID(), organization.ID(), price.ID())

	if err := c.tenantRepository.Create(tenant); err != nil {
		return err
	}

	return nil
}
