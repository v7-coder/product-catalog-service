package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/v7-coder/product-catalog-service/internal/application/product"
	"github.com/v7-coder/product-catalog-service/internal/config"
	"github.com/v7-coder/product-catalog-service/internal/infrastructure/persistence/postgres"
	"github.com/v7-coder/product-catalog-service/internal/interfaces/http/handler"
	"github.com/v7-coder/product-catalog-service/internal/interfaces/http/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	db, err := postgres.NewDB(cfg.DBUrl)
	if err != nil {
		log.Fatalf("config error: %v", err)
	}
	defer db.Close()

	log.Println("Connected to database")

	// Domain
	productRepo := postgres.NewProductRepository(db)

	// Application
	createUC := product.NewCreateProductUseCase(productRepo)

	// HTTP
	productHandler := handler.NewCreateProductHandler(createUC)

	handlers := &router.Handlers{
		Product: productHandler,
	}

	mux := http.NewServeMux()
	router.RegisterRoutes(mux, handlers)

	log.Println("Starting http server on port:", cfg.AppPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.AppPort), mux); err != nil {
		log.Fatal(err)
	}
}
