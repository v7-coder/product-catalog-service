package category

import (
	"context"
	"fmt"

	"github.com/v7-coder/product-catalog-service/internal/domain/category"
)

type UpdateCategoryCommand struct {
	ID          int
	Name        string
	Description string
}

type UpdateCategoryResult struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateCategoryUseCase struct {
	repo category.Repository
}

func NewUpdateCategoryUseCase(repo category.Repository) *UpdateCategoryUseCase {
	return &UpdateCategoryUseCase{
		repo: repo,
	}
}

func (uc *UpdateCategoryUseCase) Execute(ctx context.Context, cmd UpdateCategoryCommand) (*UpdateCategoryResult, error) {
	cat := category.NewCategory(
		cmd.Name,
		cmd.Description,
	)
	cat.ID = cmd.ID

	err := uc.repo.Update(ctx, cat)
	if err != nil {
		return nil, fmt.Errorf("update category: %w", err)
	}

	return &UpdateCategoryResult{
		ID:          cat.ID,
		Name:        cat.Name,
		Description: cat.Description,
	}, nil
}