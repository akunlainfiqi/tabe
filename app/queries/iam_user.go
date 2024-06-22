package queries

type Users struct {
	UserId   string     `json:"user_id"`
	Name     string     `json:"name"`
	Email    string     `json:"email"`
	Active   bool       `json:"active"`
	UserRole []UserRole `json:"user_role"`
}

type UserQuery interface {
	GetAll() ([]Users, error)
	GetByID(id string) (Users, error)
}
