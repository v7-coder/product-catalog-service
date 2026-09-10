package category

import (
	"context"
	"fmt"

	"github.com/v7-coder/product-catalog-service/internal/domain/category"
)

type GetCategoryCommand struct {
	ID int
}

type GetCategoryResult struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GetCategoryUseCase struct {
	repo category.Repository
}

func NewGetCategoryUseCase(repo category.Repository) *GetCategoryUseCase {
	return &GetCategoryUseCase{
		repo: repo,
	}
}

func (uc *GetCategoryUseCase) Execute(ctx context.Context, cmd GetCategoryCommand) (*GetCategoryResult, error) {
	cat, err := uc.repo.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("get category: %w", err)
	}

	return &GetCategoryResult{
		ID:          cat.ID,
		Name:        cat.Name,
		Description: cat.Description,
	}, nil
}