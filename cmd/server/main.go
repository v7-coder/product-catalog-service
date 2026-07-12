package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	appProduct "github.com/v7-coder/product-catalog-service/internal/application/product"
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

	db, err := connectWithRetry(cfg.DBUrl)
	if err != nil {
		log.Fatalf("database connection: %v", err)
	}
	defer db.Close()

	log.Println("Connected to database")

	if err := runMigrations(cfg.DBUrl); err != nil {
		log.Fatalf("migration error: %v", err)
	}

	// Domain
	productRepo, err := postgres.NewProductRepository(db)
	if err != nil {
		log.Fatalf("create product repository: %v", err)
	}

	// Application
	createUC := appProduct.NewCreateProductUseCase(productRepo)

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

func connectWithRetry(dbUrl string) (*postgres.DB, error) {
	var db *postgres.DB
	var err error

	for i := 0; i < 30; i++ {
		db, err = postgres.NewDB(dbUrl)
		if err == nil {
			log.Println("Connected to database")
			return db, nil
		}
		log.Printf("Waiting for database... (attempt %d/30)", i+1)
		time.Sleep(time.Second)
	}

	return nil, fmt.Errorf("database not available after 30 attempts: %w", err)
}

func runMigrations(dbUrl string) error {
	m, err := migrate.New("file://migrations", dbUrl)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Migrations applied successfully")

	return nil
}
