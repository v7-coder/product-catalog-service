package product

import (
	"context"
	"fmt"

	"github.com/v7-coder/product-catalog-service/internal/domain/product"
)

type UpdateProductCommand struct {
	ID          int
	Name        string
	Description string
	Price       float64
	CategoryId  int
}

type UpdateProductResult struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CategoryId  int     `json:"categoryId"`
}

type UpdateProductUseCase struct {
	repo product.Repository
}

func NewUpdateProductUseCase(repo product.Repository) *UpdateProductUseCase {
	return &UpdateProductUseCase{
		repo: repo,
	}
}

func (uc *UpdateProductUseCase) Execute(ctx context.Context, cmd UpdateProductCommand) (*UpdateProductResult, error) {
	prod := product.NewProduct(
		cmd.Name,
		cmd.Description,
		cmd.Price,
		cmd.CategoryId,
	)
	prod.ID = cmd.ID

	err := uc.repo.Update(ctx, prod)
	if err != nil {
		return nil, fmt.Errorf("update product: %w", err)
	}

	return &UpdateProductResult{
		ID:          prod.ID,
		Name:        prod.Name,
		Description: prod.Description,
		Price:       prod.Price,
		CategoryId:  prod.CategoryId,
	}, nil
}