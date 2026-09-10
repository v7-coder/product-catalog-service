package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/v7-coder/product-catalog-service/internal/domain/product"
)

type ProductRepository struct {
	db             *DB
	insertStmt     *sqlx.NamedStmt
	getByIDStmt    *sqlx.NamedStmt
	getListStmt    *sqlx.NamedStmt
	updateStmt     *sqlx.NamedStmt
	deleteStmt     *sqlx.NamedStmt
}

func NewProductRepository(db *DB) (*ProductRepository, error) {
	// Подготовленный запрос для INSERT
	insertStmt, err := db.PrepareNamed(`
		INSERT INTO products (name, description, price, category_id)
		VALUES (:name, :description, :price, :category_id)
		RETURNING id, created_at, updated_at
	`)
	if err != nil {
		return nil, fmt.Errorf("prepare insert statement: %w", err)
	}

	// Подготовленный запрос для SELECT по ID
	getByIDStmt, err := db.PrepareNamed(`
		SELECT id, name, description, price, category_id, created_at, updated_at
		FROM products
		WHERE id = :id
	`)
	if err != nil {
		return nil, fmt.Errorf("prepare get by id statement: %w", err)
	}

	// Подготовленный запрос для SELECT списка
	getListStmt, err := db.PrepareNamed(`
		SELECT id, name, description, price, category_id, created_at, updated_at
		FROM products
		ORDER BY id
		LIMIT :limit OFFSET :offset
	`)
	if err != nil {
		return nil, fmt.Errorf("prepare get list statement: %w", err)
	}

	// Подготовленный запрос для UPDATE
	updateStmt, err := db.PrepareNamed(`
		UPDATE products
		SET name = :name, description = :description, price = :price, category_id = :category_id, updated_at = NOW()
		WHERE id = :id
		RETURNING updated_at
	`)
	if err != nil {
		return nil, fmt.Errorf("prepare update statement: %w", err)
	}

	// Подготовленный запрос для DELETE
	deleteStmt, err := db.PrepareNamed(`
		DELETE FROM products
		WHERE id = :id
	`)
	if err != nil {
		return nil, fmt.Errorf("prepare delete statement: %w", err)
	}

	return &ProductRepository{
		db:             db,
		insertStmt:     insertStmt,
		getByIDStmt:    getByIDStmt,
		getListStmt:    getListStmt,
		updateStmt:     updateStmt,
		deleteStmt:     deleteStmt,
	}, nil
}

var _ product.Repository = (*ProductRepository)(nil)

func (r *ProductRepository) Save(ctx context.Context, product *product.Product) error {
	err := r.insertStmt.QueryRowxContext(ctx, product).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert product: %w", err)
	}
	return nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id int) (*product.Product, error) {
	var prod product.Product
	err := r.getByIDStmt.GetContext(ctx, &prod, map[string]interface{}{"id": id})
	if err != nil {
		return nil, fmt.Errorf("get product by id: %w", err)
	}
	return &prod, nil
}

func (r *ProductRepository) GetList(ctx context.Context, limit, offset int) ([]*product.Product, error) {
	var products []*product.Product
	err := r.getListStmt.SelectContext(ctx, &products, map[string]interface{}{"limit": limit, "offset": offset})
	if err != nil {
		return nil, fmt.Errorf("get product list: %w", err)
	}
	return products, nil
}

func (r *ProductRepository) Update(ctx context.Context, product *product.Product) error {
	err := r.updateStmt.QueryRowxContext(ctx, product).Scan(&product.UpdatedAt)
	if err != nil {
		return fmt.Errorf("update product: %w", err)
	}
	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, id int) error {
	_, err := r.deleteStmt.ExecContext(ctx, map[string]interface{}{"id": id})
	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}
	return nil
}