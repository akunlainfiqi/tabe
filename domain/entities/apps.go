package entities

type Apps struct {
	id   string
	name string
}

func NewApps(
	id,
	name string,
) *Apps {
	return &Apps{
		id:   id,
		name: name,
	}
}

func (p *Apps) ID() string {
	return p.id
}

func (p *Apps) Name() string {
	return p.name
}
