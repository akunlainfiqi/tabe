package queries

type UserRole struct {
	UserId       string `json:"user_id"`
	RoleId       string `json:"role_id"`
	BillingAccId string `json:"user_role_billing_account_id"`
}
