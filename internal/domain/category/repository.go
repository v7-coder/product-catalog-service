package category

import "context"

type Repository interface {
	Save(ctx context.Context, category *Category) error
	GetByID(ctx context.Context, id int) (*Category, error)
	GetList(ctx context.Context, limit, offset int) ([]*Category, error)
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id int) error
}