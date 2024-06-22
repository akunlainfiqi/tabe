package errors

var (
	ErrUnauthorized = NewInvariant(InvariantParam{
		Code:    4011,
		Message: "unauthorized",
	})
)
var (
	ErrInvalidBillAmount = NewInvariant(InvariantParam{
		Code:    1001,
		Message: "invalid_bill_amount",
	})
)

var (
	ErrOrganizationNotFound = NewInvariant(InvariantParam{
		Code:    4041,
		Message: "organization_not_found",
	})

	ErrAppsNotFound = NewInvariant(InvariantParam{
		Code:    4042,
		Message: "apps_not_found",
	})

	ErrBillsNotFound = NewInvariant(InvariantParam{
		Code:    4043,
		Message: "bills_not_found",
	})

	ErrTenantNotFound = NewInvariant(InvariantParam{
		Code:    4044,
		Message: "tenant_not_found",
	})

	ErrPriceNotFound = NewInvariant(InvariantParam{
		Code:    4045,
		Message: "price_not_found",
	})

	ErrTransactionNotFound = NewInvariant(InvariantParam{
		Code:    4046,
		Message: "transaction_not_found",
	})

	ErrProductNotFound = NewInvariant(InvariantParam{
		Code:    4047,
		Message: "product_not_found",
	})
)
