package queries

type App struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	ImageUrl string `json:"image_url"`
}

type AppQueries interface {
	GetAll() ([]App, error)
}
