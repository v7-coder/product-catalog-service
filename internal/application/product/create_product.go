package product

import (
	"context"
	"errors"
	"fmt"

	"github.com/gosimple/slug"
	"github.com/v7-coder/product-catalog-service/internal/domain/product"
)

type CreateProductCommand struct {
	Name        string
	Description string
	Price       float64
	CategoryId  int
}

type CreateProductResult struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CategoryId  int     `json:"categoryId"`
}

type CreateProductUseCase struct {
	repo product.Repository
}

func NewCreateProductUseCase(repo product.Repository) *CreateProductUseCase {
	return &CreateProductUseCase{
		repo: repo,
	}
}

const maxSlugAttempts = 100

func (uc *CreateProductUseCase) Execute(ctx context.Context, cmd CreateProductCommand) (*CreateProductResult, error) {
	baseSlug := slug.Make(cmd.Name)
	if baseSlug == "" {
		baseSlug = "product"
	}

	for i := 0; i < maxSlugAttempts; i++ {
		candidate := baseSlug
		if i > 0 {
			candidate = fmt.Sprintf("%s-%d", baseSlug, i+1)
		}

		prod := product.NewProduct(cmd.Name, cmd.Description, cmd.Price, cmd.CategoryId)
		prod.Slug = candidate

		err := uc.repo.Save(ctx, prod)
		if err == nil {
			return &CreateProductResult{
				ID:          prod.ID,
				Name:        prod.Name,
				Slug:        prod.Slug,
				Description: prod.Description,
				Price:       prod.Price,
				CategoryId:  prod.CategoryId,
			}, nil
		}

		if !errors.Is(err, product.ErrSlugAlreadyExists) {
			return nil, fmt.Errorf("save product: %w", err)
		}
	}

	return nil, fmt.Errorf("could not generate unique slug for %q after %d attempts", cmd.Name, maxSlugAttempts)
}
