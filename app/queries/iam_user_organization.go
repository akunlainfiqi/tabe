package queries

type IamUserOrganization struct {
	Id             string
	UserId         string
	OrganizationId string
	Level          string
}

type IamUserOrganizationQuery interface {
	IsOwner(organizationID, userID string) bool
}
