package entities

type Organization struct {
	id             string
	name           string
	identifier     string
	balance        int64
	contactName    string
	contactEmail   string
	contactPhone   string
	contactAddress string
}

func BuildOrganization(
	id,
	name,
	identifier string,
	balance int64,
	contactName,
	contactEmail,
	contactPhone,
	contactAddress string,
) *Organization {
	return &Organization{
		id:             id,
		name:           name,
		identifier:     identifier,
		balance:        balance,
		contactName:    contactName,
		contactEmail:   contactEmail,
		contactPhone:   contactPhone,
		contactAddress: contactAddress,
	}
}

func NewOrganization(
	id,
	name,
	identifier,
	contactName,
	contactEmail,
	contactPhone,
	contactAddress string,
) *Organization {
	balance := int64(0)
	return &Organization{
		id:             id,
		name:           name,
		identifier:     identifier,
		balance:        balance,
		contactName:    contactName,
		contactEmail:   contactEmail,
		contactPhone:   contactPhone,
		contactAddress: contactAddress,
	}
}

func (o *Organization) ID() string {
	return o.id
}

func (o *Organization) Name() string {
	return o.name
}

func (o *Organization) Identifier() string {
	return o.identifier
}

func (o *Organization) Balance() int64 {
	return o.balance
}

func (o *Organization) ContactName() string {
	return o.contactName
}

func (o *Organization) ContactEmail() string {
	return o.contactEmail
}

func (o *Organization) ContactPhone() string {
	return o.contactPhone
}

func (o *Organization) ContactAddress() string {
	return o.contactAddress
}

func (o *Organization) SetBalance(balance int64) {
	o.balance = balance
}
