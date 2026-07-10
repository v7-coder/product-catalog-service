package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/v7-coder/product-catalog-service/internal/domain/product"
)

type ProductRepository struct {
	db *DB
}

func NewProductRepository(db *DB) *ProductRepository {
	return &ProductRepository{db: db}
}

var _ product.Repository = (*ProductRepository)(nil)

func (r *ProductRepository) Save(ctx context.Context, product *product.Product) error {
	query, args, err := sq.Insert("products").
		Columns("name", "description", "price", "category_id").
		Values(product.Name, product.Description, product.Price, product.CategoryId).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("build insert query: %w", err)
	}

	err = r.db.QueryRowContext(ctx, query, args...).Scan(&product.ID)
	if err != nil {
		return fmt.Errorf("insert product: %w", err)
	}

	return nil
}
