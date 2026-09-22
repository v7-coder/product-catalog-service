package product

type Product struct {
	ID          int     `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	Description string  `json:"description" db:"description"`
	Price       float64 `json:"price" db:"price"`
	CategoryId  int     `json:"categoryId" db:"category_id"`
}

func NewProduct(name, description string, price float64, categoryId int) *Product {
	return &Product{
		Name:        name,
		Description: description,
		Price:       price,
		CategoryId:  categoryId,
	}
}
