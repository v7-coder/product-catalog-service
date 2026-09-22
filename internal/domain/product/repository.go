package product

import "context"

type Repository interface {
	Save(ctx context.Context, product *Product) error
}
