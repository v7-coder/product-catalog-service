package category

import (
	"context"
	"fmt"

	"github.com/v7-coder/product-catalog-service/internal/domain/category"
)

type CreateCategoryCommand struct {
	Name        string
	Description string
}

type CreateCategoryResult struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateCategoryUseCase struct {
	repo category.Repository
}

func NewCreateCategoryUseCase(repo category.Repository) *CreateCategoryUseCase {
	return &CreateCategoryUseCase{
		repo: repo,
	}
}

func (uc *CreateCategoryUseCase) Execute(ctx context.Context, cmd CreateCategoryCommand) (*CreateCategoryResult, error) {
	cat := category.NewCategory(
		cmd.Name,
		cmd.Description,
	)

	err := uc.repo.Save(ctx, cat)
	if err != nil {
		return nil, fmt.Errorf("save category: %w", err)
	}

	return &CreateCategoryResult{
		ID:          cat.ID,
		Name:        cat.Name,
		Description: cat.Description,
	}, nil
}