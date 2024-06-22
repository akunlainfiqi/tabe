package entities

var (
	RoleAdmin = Role{
		id:   "1",
		name: "Admin",
	}
	RoleBillingAccountOwner = Role{
		id:   "2",
		name: "Billing Account Owner",
	}
	RoleBillingAccountStaff = Role{
		id:   "3",
		name: "Billing Account Staff",
	}
	RoleBillingAccountUser = Role{
		id:   "4",
		name: "Billing Account User",
	}
)

type Role struct {
	id   string
	name string
}

func NewRole(id, name string) *Role {
	return &Role{
		id:   id,
		name: name,
	}
}

func (r *Role) ID() string {
	return r.id
}

func (r *Role) Name() string {
	return r.name
}
