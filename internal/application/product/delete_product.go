package product

import (
	"context"
	"fmt"

	"github.com/v7-coder/product-catalog-service/internal/domain/product"
)

type DeleteProductCommand struct {
	ID int
}

type DeleteProductResult struct {
	Success bool `json:"success"`
}

type DeleteProductUseCase struct {
	repo product.Repository
}

func NewDeleteProductUseCase(repo product.Repository) *DeleteProductUseCase {
	return &DeleteProductUseCase{
		repo: repo,
	}
}

func (uc *DeleteProductUseCase) Execute(ctx context.Context, cmd DeleteProductCommand) (*DeleteProductResult, error) {
	err := uc.repo.Delete(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("delete product: %w", err)
	}

	return &DeleteProductResult{
		Success: true,
	}, nil
}