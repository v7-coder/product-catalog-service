package product

import "errors"

var ErrSlugAlreadyExists = errors.New("product with this slug already exists")
