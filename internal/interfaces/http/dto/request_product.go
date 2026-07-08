package dto

type CreateProductRequest struct {
	Name string `json:"name" validate:"required,min=3,max=255"`
}
