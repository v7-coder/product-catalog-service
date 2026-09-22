package router

import (
	"net/http"

	"github.com/v7-coder/product-catalog-service/internal/interfaces/http/handler"
)

type Handlers struct {
	Product *handler.ProductHandler
}

func RegisterRoutes(mux *http.ServeMux, handlers *Handlers) {
	// POST /api/v1/products - создание продукта
	mux.HandleFunc("POST /api/v1/products", handlers.Product.Create)
}
