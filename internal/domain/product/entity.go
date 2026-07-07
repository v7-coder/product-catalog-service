package product

type Product struct {
	ID   int    `json:"id"`
	Name string `json:"name" validate:"reqquired, min=3,max=255"`
}
