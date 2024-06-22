package queries

type Price struct {
	ID         string `json:"id"`
	ProductId  string `json:"product_id"`
	Price      int64  `json:"price"`
	Reccurence string `json:"reccurence"`
}

type PriceQuery interface {
	FindAll() ([]Price, error)
}
