package queries

type Organization struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
	Balance    int64  `json:"balance"`
}

type OrganizationQuery interface {
	GetByID(id string) (Organization, error)
}
