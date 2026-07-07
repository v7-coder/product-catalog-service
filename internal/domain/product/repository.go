package product

import "context"

type Repository interface {
	Save(context.Context) error
}