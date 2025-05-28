package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/skiba/lamp_store/orders_service/internal/handler"
	corsmiddleware "github.com/skiba/lamp_store/orders_service/internal/middleware"
	"github.com/skiba/lamp_store/orders_service/internal/repository/postgres"
	"github.com/skiba/lamp_store/orders_service/internal/service"
	"github.com/skiba/lamp_store/orders_service/pkg/config"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.Database.GetDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Автоматически применяем миграции
	m, err := migrate.New(
		"file:///app/migrations",
		cfg.Database.GetURLDSN(),
	)
	if err != nil {
		log.Printf("Failed to init migrate: %v", err)
	} else {
		err = m.Up()
		if err != nil {
			if err.Error() == "no change" {
				// Ничего не делаем, это нормально
			} else if err.Error() == "Dirty database version 2. Fix and force version." {
				// Исправляем "грязное" состояние
				if err := m.Force(2); err != nil {
					log.Printf("Failed to force version: %v", err)
				}
			} else {
				log.Printf("Failed to apply migrations: %v", err)
			}
		}
	}

	orderRepo := postgres.NewOrderRepository(db)
	productsClient := service.NewProductsClient("http://products_service:8000")
	orderService := service.NewOrderService(orderRepo, productsClient)
	orderHandler := handler.NewOrderHandler(orderService)

	r := chi.NewRouter()

	// Middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.RequestID)
	r.Use(corsmiddleware.Cors)

	orderHandler.Register(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
	}

	log.Printf("Server is starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
