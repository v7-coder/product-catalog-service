package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/v7-coder/product-catalog-service/internal/domain/product"
)

type ProductRepository struct {
	db         *DB
	insertStmt *sqlx.NamedStmt
}

func NewProductRepository(db *DB) (*ProductRepository, error) {
	// Подготовленный запрос для INSERT
	insertStmt, err := db.PrepareNamed(`
		INSERT INTO products (name, slug, description, price, category_id)
		VALUES (:name, :slug, :description, :price, :category_id)
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

func (r *ProductRepository) Save(ctx context.Context, p *product.Product) error {
	err := r.insertStmt.QueryRowxContext(ctx, p).Scan(&p.ID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return product.ErrSlugAlreadyExists
		}

		return fmt.Errorf("insert product: %w", err)
	}
	return nil
}
