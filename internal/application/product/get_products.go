package product

import (
	"context"
	"fmt"

	"github.com/v7-coder/product-catalog-service/internal/domain/product"
)

type GetProductsCommand struct {
	Limit  int
	Offset int
}

type GetProductsResult struct {
	Products []*product.Product `json:"products"`
}

type GetProductsUseCase struct {
	repo product.Repository
}

func NewGetProductsUseCase(repo product.Repository) *GetProductsUseCase {
	return &GetProductsUseCase{
		repo: repo,
	}
}

func (uc *GetProductsUseCase) Execute(ctx context.Context, cmd GetProductsCommand) (*GetProductsResult, error) {
	products, err := uc.repo.GetList(ctx, cmd.Limit, cmd.Offset)
	if err != nil {
		return nil, fmt.Errorf("get products: %w", err)
	}

	return &GetProductsResult{
		Products: products,
	}, nil
}