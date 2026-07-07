package product

import (
	"errors"

	"github.com/v7-coder/product-catalog-service/internal/pkg/validator"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ValidateProduct(p *Product) error {
	if p == nil {
		return errors.New("product is nil")
	}

	return validator.Validate(p)
}
