package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/v7-coder/product-catalog-service/internal/domain/category"
)

type CategoryRepository struct {
	db             *DB
	insertStmt     *sqlx.NamedStmt
	getByIDStmt    *sqlx.NamedStmt
	getListStmt    *sqlx.NamedStmt
	updateStmt     *sqlx.NamedStmt
	deleteStmt     *sqlx.NamedStmt
}

func NewCategoryRepository(db *DB) (*CategoryRepository, error) {
	// Подготовленный запрос для INSERT
	insertStmt, err := db.PrepareNamed(`
		INSERT INTO categories (name, description)
		VALUES (:name, :description)
		RETURNING id, created_at, updated_at
	`)
	if err != nil {
		return nil, fmt.Errorf("prepare insert statement: %w", err)
	}

	// Подготовленный запрос для SELECT по ID
	getByIDStmt, err := db.PrepareNamed(`
		SELECT id, name, description, created_at, updated_at
		FROM categories
		WHERE id = :id
	`)
	if err != nil {
		return nil, fmt.Errorf("prepare get by id statement: %w", err)
	}

	// Подготовленный запрос для SELECT списка
	getListStmt, err := db.PrepareNamed(`
		SELECT id, name, description, created_at, updated_at
		FROM categories
		ORDER BY id
		LIMIT :limit OFFSET :offset
	`)
	if err != nil {
		return nil, fmt.Errorf("prepare get list statement: %w", err)
	}

	// Подготовленный запрос для UPDATE
	updateStmt, err := db.PrepareNamed(`
		UPDATE categories
		SET name = :name, description = :description, updated_at = NOW()
		WHERE id = :id
		RETURNING updated_at
	`)
	if err != nil {
		return nil, fmt.Errorf("prepare update statement: %w", err)
	}

	// Подготовленный запрос для DELETE
	deleteStmt, err := db.PrepareNamed(`
		DELETE FROM categories
		WHERE id = :id
	`)
	if err != nil {
		return nil, fmt.Errorf("prepare delete statement: %w", err)
	}

	return &CategoryRepository{
		db:             db,
		insertStmt:     insertStmt,
		getByIDStmt:    getByIDStmt,
		getListStmt:    getListStmt,
		updateStmt:     updateStmt,
		deleteStmt:     deleteStmt,
	}, nil
}

var _ category.Repository = (*CategoryRepository)(nil)

func (r *CategoryRepository) Save(ctx context.Context, category *category.Category) error {
	err := r.insertStmt.QueryRowxContext(ctx, category).Scan(&category.ID, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert category: %w", err)
	}
	return nil
}

func (r *CategoryRepository) GetByID(ctx context.Context, id int) (*category.Category, error) {
	var cat category.Category
	err := r.getByIDStmt.GetContext(ctx, &cat, map[string]interface{}{"id": id})
	if err != nil {
		return nil, fmt.Errorf("get category by id: %w", err)
	}
	return &cat, nil
}

func (r *CategoryRepository) GetList(ctx context.Context, limit, offset int) ([]*category.Category, error) {
	var categories []*category.Category
	err := r.getListStmt.SelectContext(ctx, &categories, map[string]interface{}{"limit": limit, "offset": offset})
	if err != nil {
		return nil, fmt.Errorf("get category list: %w", err)
	}
	return categories, nil
}

func (r *CategoryRepository) Update(ctx context.Context, category *category.Category) error {
	err := r.updateStmt.QueryRowxContext(ctx, category).Scan(&category.UpdatedAt)
	if err != nil {
		return fmt.Errorf("update category: %w", err)
	}
	return nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id int) error {
	_, err := r.deleteStmt.ExecContext(ctx, map[string]interface{}{"id": id})
	if err != nil {
		return fmt.Errorf("delete category: %w", err)
	}
	return nil
}