package queries

type Product struct {
	ID        string  `json:"id"`
	AppId     string  `json:"app_id"`
	Name      string  `json:"name"`
	TierName  string  `json:"tier_name"`
	TierIndex int     `json:"tier_index"`
	Price     []Price `json:"price"`
}

type ProductQuery interface {
	FindAll() ([]Product, error)
	FindByID(id string) (Product, error)
	FindByAppID(appID string) ([]Product, error)
}
