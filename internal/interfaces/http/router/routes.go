package router

import (
	"net/http"

	"github.com/v7-coder/product-catalog-service/internal/interfaces/http/handler"
)

type Handlers struct {
	Product  *handler.ProductHandler
	Category *handler.CategoryHandler
}

func RegisterRoutes(mux *http.ServeMux, handlers *Handlers) {
	// POST /api/v1/products - создание продукта
	mux.HandleFunc("POST /api/v1/products", handlers.Product.Create)

	// GET /api/v1/products/{productId} - получение продукта по ID
	mux.HandleFunc("GET /api/v1/products/{productId}", handlers.Product.Get)

	// GET /api/v1/products - список продуктов с пагинацией
	mux.HandleFunc("GET /api/v1/products", handlers.Product.GetList)

	// PUT /api/v1/products/{productId} - обновление продукта
	mux.HandleFunc("PUT /api/v1/products/{productId}", handlers.Product.Update)

	// DELETE /api/v1/products/{productId} - удаление продукта
	mux.HandleFunc("DELETE /api/v1/products/{productId}", handlers.Product.Delete)

	// POST /api/v1/categories - создание категории
	mux.HandleFunc("POST /api/v1/categories", handlers.Category.Create)

	// GET /api/v1/categories/{categoryId} - получение категории по ID
	mux.HandleFunc("GET /api/v1/categories/{categoryId}", handlers.Category.Get)

	// GET /api/v1/categories - список категорий с пагинацией
	mux.HandleFunc("GET /api/v1/categories", handlers.Category.GetList)

	// PUT /api/v1/categories/{categoryId} - обновление категории
	mux.HandleFunc("PUT /api/v1/categories/{categoryId}", handlers.Category.Update)

	// DELETE /api/v1/categories/{categoryId} - удаление категории
	mux.HandleFunc("DELETE /api/v1/categories/{categoryId}", handlers.Category.Delete)
}
