package entities

type Tenant struct {
	id             string
	name           string
	appID          string
	organizationId string
	activeUntil    int64
	priceId        string
}

func BuildTenant(
	id,
	name,
	appID,
	organizationId,
	priceId string,
	activeUntil int64,
) *Tenant {
	return &Tenant{
		id:             id,
		name:           name,
		appID:          appID,
		organizationId: organizationId,
		activeUntil:    activeUntil,
		priceId:        priceId,
	}
}

func NewTenant(
	id,
	name,
	appID,
	organizationId,
	priceId string,
) *Tenant {
	return &Tenant{
		id:             id,
		name:           name,
		appID:          appID,
		organizationId: organizationId,
		activeUntil:    0,
		priceId:        priceId,
	}
}

func (t *Tenant) ID() string {
	return t.id
}

func (t *Tenant) Name() string {
	return t.name
}

func (t *Tenant) AppID() string {
	return t.appID
}

func (t *Tenant) OrganizationID() string {
	return t.organizationId
}

func (t *Tenant) ActiveUntil() int64 {
	return t.activeUntil
}

func (t *Tenant) SetActiveUntil(activeUntil int64) {
	t.activeUntil = activeUntil
}

func (t *Tenant) PriceID() string {
	return t.priceId
}
