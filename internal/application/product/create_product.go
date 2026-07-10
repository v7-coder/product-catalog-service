package product

import (
	"context"
	"fmt"

	"github.com/v7-coder/product-catalog-service/internal/domain/product"
	"github.com/v7-coder/product-catalog-service/internal/pkg/validator"
)

type CreateProductCommand struct {
	Name        string  `validate:"required,min=3,max=255"`
	Description string  `json:"description" validate:"min=3,max=1000"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	CategoryId  int     `json:"categoryId" validate:"required,gt=0"`
}

type CreateProductResult struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CategoryId  int     `json:"categoryId"`
}

type CreateProductUseCase struct {
	repo product.Repository
}

func NewCreateProductUseCase(repo product.Repository) *CreateProductUseCase {
	return &CreateProductUseCase{repo: repo}
}

func (uc *CreateProductUseCase) Execute(ctx context.Context, cmd CreateProductCommand) (*CreateProductResult, error) {
	if err := validator.Validate(cmd); err != nil {
		return nil, fmt.Errorf("validation CreateProductCommand failed: %w", err)
	}

	product := product.NewProduct(cmd.Name, cmd.Description, cmd.Price, cmd.CategoryId)

	err := uc.repo.Save(ctx, product)
	if err != nil {
		return nil, err
	}

	return &CreateProductResult{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CategoryId:  product.CategoryId,
	}, nil
}
