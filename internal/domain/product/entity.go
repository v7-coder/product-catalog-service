package product

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required,min=3,max=255"`
	Description string  `json:"description" validate:"min=3,max=1000"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	CategoryId  int     `json:"categoryId" validate:"required,gt=0"`
}

func NewProduct(name, description string, price float64, categoryId int) *Product {
	return &Product{
		Name:        name,
		Description: description,
		Price:       price,
		CategoryId:  categoryId,
	}
}
