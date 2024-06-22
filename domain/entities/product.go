package entities

type Product struct {
	id        string
	app       Apps
	tierName  string
	tierIndex int
}

func NewProduct(
	id string,
	app Apps,
	tierName string,
	tierIndex int,
) *Product {
	return &Product{
		id:        id,
		app:       app,
		tierName:  tierName,
		tierIndex: tierIndex,
	}
}

func (p *Product) ID() string {
	return p.id
}

func (p *Product) App() *Apps {
	return &p.app
}

func (p *Product) TierName() string {
	return p.tierName
}

func (p *Product) TierIndex() int {
	return p.tierIndex
}
