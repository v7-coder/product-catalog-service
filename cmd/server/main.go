package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
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
	runtime.GOMAXPROCS(runtime.NumCPU())

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	appConfig, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", slog.Any("error", err))
		os.Exit(1)
	}

	if err := run(logger, appConfig); err != nil {
		logger.Error("application failed", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("application stopped gracefully")
}

func run(logger *slog.Logger, appConfig *config.Config) error {
	db, err := connectWithRetry(logger, appConfig.DBUrl)
	if err != nil {
		return fmt.Errorf("database connection: %w", err)
	}
	defer db.Close()

	if err := runMigrations(logger, appConfig.DBUrl); err != nil {
		return fmt.Errorf("migration error: %w", err)
	}

	server, err := configureServer(appConfig, db)
	if err != nil {
		return fmt.Errorf("configure server: %w", err)
	}

	errChan := make(chan error, 1)

	go func() {
		logger.Info("server starting", slog.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("server error: %w", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		logger.Info("received shutdown signal", slog.String("signal", sig.String()))
	case err := <-errChan:
		return fmt.Errorf("server failed: %w", err)
	}

	logger.Info("shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown error: %w", err)
	}

	logger.Info("server stopped")
	return nil
}

func connectWithRetry(logger *slog.Logger, dbUrl string) (*postgres.DB, error) {
	var db *postgres.DB
	var err error

	for i := 0; i < 30; i++ {
		db, err = postgres.NewDB(dbUrl)
		if err == nil {
			logger.Info("connected to database")
			return db, nil
		}
		logger.Warn("waiting for database...",
			slog.Int("attempt", i+1),
			slog.Int("max_attempts", 30),
		)
		time.Sleep(time.Second)
	}

	return nil, fmt.Errorf("database not available after 30 attempts: %w", err)
}

func configureServer(appConfig *config.Config, db *postgres.DB) (*http.Server, error) {
	productRepo, err := postgres.NewProductRepository(db)
	if err != nil {
		return nil, fmt.Errorf("create product repository: %w", err)
	}

	createUC := appProduct.NewCreateProductUseCase(productRepo)

	productHandler := handler.NewCreateProductHandler(createUC)

	handlers := &router.Handlers{
		Product: productHandler,
	}

	mux := http.NewServeMux()
	router.RegisterRoutes(mux, handlers)

	return &http.Server{
		Addr:         fmt.Sprintf(":%s", appConfig.AppPort),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}, nil
}

func runMigrations(logger *slog.Logger, dbUrl string) error {
	m, err := migrate.New("file://migrations", dbUrl)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	logger.Info("migrations applied successfully")
	return nil
}
