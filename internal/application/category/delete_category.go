package category

import (
	"context"
	"fmt"

	"github.com/v7-coder/product-catalog-service/internal/domain/category"
)

type DeleteCategoryCommand struct {
	ID int
}

type DeleteCategoryResult struct {
	Success bool `json:"success"`
}

type DeleteCategoryUseCase struct {
	repo category.Repository
}

func NewDeleteCategoryUseCase(repo category.Repository) *DeleteCategoryUseCase {
	return &DeleteCategoryUseCase{
		repo: repo,
	}
}

func (uc *DeleteCategoryUseCase) Execute(ctx context.Context, cmd DeleteCategoryCommand) (*DeleteCategoryResult, error) {
	err := uc.repo.Delete(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("delete category: %w", err)
	}

	return &DeleteCategoryResult{
		Success: true,
	}, nil
}