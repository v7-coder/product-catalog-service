package category

import (
	"context"
	"fmt"

	"github.com/v7-coder/product-catalog-service/internal/domain/category"
)

type GetCategoriesCommand struct {
	Limit  int
	Offset int
}

type GetCategoriesResult struct {
	Categories []*category.Category `json:"categories"`
}

type GetCategoriesUseCase struct {
	repo category.Repository
}

func NewGetCategoriesUseCase(repo category.Repository) *GetCategoriesUseCase {
	return &GetCategoriesUseCase{
		repo: repo,
	}
}

func (uc *GetCategoriesUseCase) Execute(ctx context.Context, cmd GetCategoriesCommand) (*GetCategoriesResult, error) {
	categories, err := uc.repo.GetList(ctx, cmd.Limit, cmd.Offset)
	if err != nil {
		return nil, fmt.Errorf("get categories: %w", err)
	}

	return &GetCategoriesResult{
		Categories: categories,
	}, nil
}