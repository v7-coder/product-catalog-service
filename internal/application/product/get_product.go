package product

import (
	"context"
	"fmt"

	"github.com/v7-coder/product-catalog-service/internal/domain/product"
)

type GetProductCommand struct {
	ID int
}

type GetProductResult struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CategoryId  int     `json:"categoryId"`
}

type GetProductUseCase struct {
	repo product.Repository
}

func NewGetProductUseCase(repo product.Repository) *GetProductUseCase {
	return &GetProductUseCase{
		repo: repo,
	}
}

func (uc *GetProductUseCase) Execute(ctx context.Context, cmd GetProductCommand) (*GetProductResult, error) {
	prod, err := uc.repo.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("get product: %w", err)
	}

	return &GetProductResult{
		ID:          prod.ID,
		Name:        prod.Name,
		Description: prod.Description,
		Price:       prod.Price,
		CategoryId:  prod.CategoryId,
	}, nil
}