package commands

import (
	"saas-billing/app/services"
	"saas-billing/domain/entities"
	"saas-billing/domain/repositories"
	"saas-billing/errors"
	"time"

	"github.com/google/uuid"
)

type CreateTenantOnboardingRequest struct {
	priceId        string
	organizationId string
	tenantId       string
	tenantName     string
	userId         string
}

func NewCreateTenantOnboardingRequest(
	priceId,
	organizationId,
	tenantId,
	tenantName,
	userId string,
) (*CreateTenantOnboardingRequest, error) {
	return &CreateTenantOnboardingRequest{
		priceId:        priceId,
		organizationId: organizationId,
		tenantId:       tenantId,
		tenantName:     tenantName,
		userId:         userId,
	}, nil
}

type CreateTenantOnboardingCommand struct {
	tenantRepository       repositories.TenantRepository
	organizationRepository repositories.OrganizationRepository
	priceRepository        repositories.PriceRepository
	billsRepository        repositories.BillsRepository

	iamOrganizationRepository repositories.IamOrganizationRepository

	midtransService services.Midtrans
}

func NewCreateTenantOnboardingCommand(
	tenantRepository repositories.TenantRepository,
	organizationRepository repositories.OrganizationRepository,
	priceRepository repositories.PriceRepository,
	billsRepository repositories.BillsRepository,
	iamOrganizationRepository repositories.IamOrganizationRepository,
	midtransService services.Midtrans,
) *CreateTenantOnboardingCommand {
	return &CreateTenantOnboardingCommand{
		tenantRepository:          tenantRepository,
		organizationRepository:    organizationRepository,
		priceRepository:           priceRepository,
		billsRepository:           billsRepository,
		iamOrganizationRepository: iamOrganizationRepository,
		midtransService:           midtransService,
	}
}

func (c *CreateTenantOnboardingCommand) Execute(req *CreateTenantOnboardingRequest) (interface{}, error) {
	isOrganizationOwner := c.iamOrganizationRepository.IsOwner(req.organizationId, req.userId)
	if !isOrganizationOwner {
		return nil, errors.ErrUnauthorized
	}

	organization, err := c.iamOrganizationRepository.GetByID(req.organizationId)
	if err != nil {
		return nil, err
	}

	price, err := c.priceRepository.GetByID(req.priceId)
	if err != nil {
		return nil, err
	}

	product := price.Product()
	
	tenant := entities.NewTenant(req.tenantId, req.tenantName, product.ID(), organization.ID, price.ID())

	if err := c.tenantRepository.Create(tenant); err != nil {
		return nil, err
	}

	billId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	bill, err := entities.NewBills(
		billId.String(),
		organization.ID,
		tenant.ID(),
		price.Price(),
		0,
		time.Now().AddDate(0, 0, 10).Unix(),
	)
	if err != nil {
		return nil, err
	}

	if err := c.billsRepository.Create(bill); err != nil {
		return nil, err
	}

	res, err := c.midtransService.CreateTransaction(bill.ID(), bill.Total())
	if err != nil {
		return nil, err
	}

	return res, nil
}
