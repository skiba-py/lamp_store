package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/skiba/lamp_store/products_service/internal/handler"
	"github.com/skiba/lamp_store/products_service/internal/repository/postgres"
	"github.com/skiba/lamp_store/products_service/internal/service"
	"github.com/skiba/lamp_store/products_service/pkg/config"
	"github.com/skiba/lamp_store/products_service/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// Config
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Logger
	logger, err := logger.NewLogger(cfg.Logging.Level, cfg.Logging.Format)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// db connect
	db, err := sql.Open("postgres", cfg.Database.GetDSN())
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// ping db
	if err := db.Ping(); err != nil {
		logger.Fatal("Failed to ping database", zap.Error(err))
	}

	repo := postgres.NewProductRepository(db)

	svc := service.NewProductService(repo)

	productHandler := handler.NewProductHandler(svc)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)

	r.Route("/api/v1/products", func(r chi.Router) {
		r.Post("/", productHandler.CreateProduct)
		r.Get("/", productHandler.ListProducts)
		r.Get("/{id}", productHandler.GetProduct)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})

	addr := cfg.Server.GetServerAddr()
	logger.Info("Starting server", zap.String("address", addr))
	if err := http.ListenAndServe(addr, r); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
