package entities

type Price struct {
	id         string
	product    Product
	price      float64
	reccurence string
}

const (
	ProductRecurrenceMonthly = "monthly"
	ProductRecurrenceYearly  = "yearly"
)

func NewPrice(
	id string,
	product Product,
	price float64,
	reccurence string,
) *Price {
	return &Price{
		id:         id,
		product:    product,
		price:      price,
		reccurence: reccurence,
	}
}

func (p *Price) ID() string {
	return p.id
}

func (p *Price) Product() Product {
	return p.product
}

func (p *Price) Price() float64 {
	return p.price
}

func (p *Price) Recurrence() string {
	return p.reccurence
}
