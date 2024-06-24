package queries

type App struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type AppQueries interface {
	GetAll() ([]App, error)
}
