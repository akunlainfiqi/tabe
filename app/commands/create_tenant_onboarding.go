package commands

import (
	"encoding/json"
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

	midtransService  services.Midtrans
	publisherService services.Publisher
}

func NewCreateTenantOnboardingCommand(
	tenantRepository repositories.TenantRepository,
	organizationRepository repositories.OrganizationRepository,
	priceRepository repositories.PriceRepository,
	billsRepository repositories.BillsRepository,
	iamOrganizationRepository repositories.IamOrganizationRepository,
	midtransService services.Midtrans,
	publisherService services.Publisher,
) *CreateTenantOnboardingCommand {
	return &CreateTenantOnboardingCommand{
		tenantRepository:          tenantRepository,
		organizationRepository:    organizationRepository,
		priceRepository:           priceRepository,
		billsRepository:           billsRepository,
		iamOrganizationRepository: iamOrganizationRepository,
		midtransService:           midtransService,
		publisherService:          publisherService,
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

	org, err := c.organizationRepository.FindByID(organization.ID)
	if err != nil {
		return nil, err
	}

	if org == nil {
		org := entities.NewOrganization(organization.ID, organization.Name, organization.Identifier, "", "", "", "")

		if err := c.organizationRepository.Create(org); err != nil {
			return nil, err
		}
	}

	price, err := c.priceRepository.GetByID(req.priceId)
	if err != nil {
		return nil, err
	}

	product := price.Product()

	billId, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	tenant := entities.NewTenant(req.tenantId, req.tenantName, product.App().ID(), organization.ID, price.ID())

	if err := c.tenantRepository.Create(tenant); err != nil {
		return nil, err
	}

	if org.Balance() > price.Price() {
		balanceUsed := price.Price()

		org.SetBalance(org.Balance() - balanceUsed)

		bill, err := entities.NewBills(
			billId.String(),
			organization.ID,
			tenant.ID(),
			price.ID(),
			price.Price(),
			balanceUsed,
			time.Now().Add(1*time.Hour).Unix(),
			string(entities.BillTypeNewSubscription),
		)

		if err != nil {
			return nil, err
		}

		if err := c.billsRepository.Create(bill); err != nil {
			return nil, err
		}

		if err := c.organizationRepository.Update(org); err != nil {
			return nil, err
		}

		if price.Recurrence() == entities.ProductRecurrenceMonthly {
			tenant.SetActiveUntil(time.Now().AddDate(0, 1, 0).Unix())
		} else if price.Recurrence() == entities.ProductRecurrenceYearly {
			tenant.SetActiveUntil(time.Now().AddDate(1, 0, 0).Unix())
		}

		if err := c.tenantRepository.Create(tenant); err != nil {
			return nil, err
		}

		pl := services.TenantPaidPayload{
			TenantID:  tenant.ID(),
			ProductID: product.ID(),
			Timestamp: time.Now(),
		}

		payload, err := json.Marshal(pl)
		if err != nil {
			return nil, err
		}

		if err := c.publisherService.Publish("billing_paid", payload); err != nil {
			return nil, err
		}

		return nil, nil
	}

	bill, err := entities.NewBills(
		billId.String(),
		organization.ID,
		tenant.ID(),
		price.ID(),
		price.Price(),
		0,
		time.Now().Add(1*time.Hour).Unix(),
		string(entities.BillTypeNewSubscription),
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
