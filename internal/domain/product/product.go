package product

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CategoryId  int     `json:"categoryId"`
}

func NewProduct(name, description string, price float64, categoryId int) *Product {
	return &Product{
		Name:        name,
		Description: description,
		Price:       price,
		CategoryId:  categoryId,
	}
}
