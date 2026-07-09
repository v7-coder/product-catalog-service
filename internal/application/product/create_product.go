package product

import (
	"context"
	"fmt"

	"github.com/v7-coder/product-catalog-service/internal/domain/product"
	"github.com/v7-coder/product-catalog-service/internal/pkg/validator"
)

type CreateProductCommand struct {
	Name string `validate:"required,min=3"`
}

type CreateProductResult struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
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

	product := product.NewProduct(cmd.Name)

	err := uc.repo.Save(ctx, product)
	if err != nil {
		return nil, err
	}

	return &CreateProductResult{
		ID:   product.ID,
		Name: product.Name,
	}, nil
}
