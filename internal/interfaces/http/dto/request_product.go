package dto

type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=3,max=255"`
	Description string  `json:"description" validate:"omitempty,min=3,max=1000"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	CategoryId  int     `json:"categoryId" validate:"required,gt=0"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name" validate:"omitempty,min=3,max=255"`
	Description string  `json:"description" validate:"omitempty,min=3,max=1000"`
	Price       float64 `json:"price" validate:"omitempty,gt=0"`
	CategoryId  int     `json:"categoryId" validate:"omitempty,gt=0"`
}