package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/v7-coder/product-catalog-service/internal/domain/product"
)

type ProductRepository struct {
	db         *DB
	insertStmt *sqlx.NamedStmt
}

func NewProductRepository(db *DB) (*ProductRepository, error) {
	// Подготовленный запрос для INSERT
	insertStmt, err := db.PrepareNamed(`
		INSERT INTO products (name, description, price, category_id)
		VALUES (:name, :description, :price, :category_id)
		RETURNING id
	`)
	if err != nil {
		return nil, fmt.Errorf("prepare insert statement: %w", err)
	}

	return &ProductRepository{
		db:         db,
		insertStmt: insertStmt,
	}, nil
}

var _ product.Repository = (*ProductRepository)(nil)

func (r *ProductRepository) Save(ctx context.Context, product *product.Product) error {
	err := r.insertStmt.QueryRowxContext(ctx, product).Scan(&product.ID)
	if err != nil {
		return fmt.Errorf("insert product: %w", err)
	}
	return nil
}
