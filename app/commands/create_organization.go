package commands

import (
	"saas-billing/domain/entities"
	"saas-billing/domain/repositories"
)

type CreateOrganizationRequest struct {
	orgId          string
	contactName    string
	contactEmail   string
	contactPhone   string
	contactAddress string
}

func NewCreateOrganizationRequest(
	orgId,
	contactName,
	contactEmail,
	contactPhone,
	contactAddress string,
) *CreateOrganizationRequest {
	return &CreateOrganizationRequest{
		orgId:          orgId,
		contactName:    contactName,
		contactEmail:   contactEmail,
		contactPhone:   contactPhone,
		contactAddress: contactAddress,
	}
}

type CreateOrganizationCommand struct {
	organizationRepository    repositories.OrganizationRepository
	iamOrganizationRepository repositories.IamOrganizationRepository
}

func NewCreateOrganizationCommand(
	organizationRepository repositories.OrganizationRepository,
	iamOrganizationRepository repositories.IamOrganizationRepository,
) *CreateOrganizationCommand {
	return &CreateOrganizationCommand{
		organizationRepository:    organizationRepository,
		iamOrganizationRepository: iamOrganizationRepository,
	}
}

func (c *CreateOrganizationCommand) Execute(req *CreateOrganizationRequest) error {
	iamOrganization, err := c.iamOrganizationRepository.GetByID(req.orgId)
	if err != nil {
		return err
	}
	organization := entities.NewOrganization(
		req.orgId,
		iamOrganization.Name,
		iamOrganization.Identifier,
		req.contactName,
		req.contactEmail,
		req.contactPhone,
		req.contactAddress,
	)

	if err := c.organizationRepository.Create(organization); err != nil {
		return err
	}

	return nil
}
