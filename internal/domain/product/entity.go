package product

type Product struct {
	ID   int    `json:"id"`
	Name string `json:"name" validate:"required, min=3,max=255"`
}

func NewProduct(name string) *Product {
	return &Product{}
}
