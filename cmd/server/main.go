package main

import (
	"log"
	"net/http"

	"github.com/v7-coder/product-catalog-service/internal/application/product"
	"github.com/v7-coder/product-catalog-service/internal/interfaces/http/handler"
	"github.com/v7-coder/product-catalog-service/internal/interfaces/http/router"
)

func main() {
	createUC := product.NewCreateProductUseCase()

	productHandler := handler.NewCreateProductHandler(createUC)

	handlers := &router.Handlers{
		Product: productHandler,
	}

	mux := http.NewServeMux()
	router.RegisterRoutes(mux, handlers)

	log.Println("Server starting on :8080") // Println вместо Panicln
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
