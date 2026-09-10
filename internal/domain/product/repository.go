package product

import "context"

type Repository interface {
	Save(ctx context.Context, product *Product) error
	GetByID(ctx context.Context, id int) (*Product, error)
	GetList(ctx context.Context, limit, offset int) ([]*Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id int) error
}