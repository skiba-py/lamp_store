package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/skiba/lamp_store/products_service/internal/handler"
	corsmiddleware "github.com/skiba/lamp_store/products_service/internal/middleware"
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

	// Автоматически применяем миграции
	m, err := migrate.New(
		"file:///app/migrations",
		cfg.Database.GetURLDSN(),
	)
	if err != nil {
		logger.Error("Failed to init migrate", zap.Error(err))
	} else {
		err = m.Up()
		if err != nil {
			if err.Error() == "no change" {
				// Ничего не делаем, это нормально
			} else if err.Error() == "Dirty database version 2. Fix and force version." {
				// Исправляем "грязное" состояние
				if err := m.Force(2); err != nil {
					logger.Error("Failed to force version", zap.Error(err))
				}
			} else {
				logger.Error("Failed to apply migrations", zap.Error(err))
			}
		}
	}

	repo := postgres.NewProductRepository(db)

	svc := service.NewProductService(repo)

	productHandler := handler.NewProductHandler(svc)

	// Создаём сервис и хендлер для резервирования
	reservationRepo := postgres.NewReservationRepository(db)
	reservationService := service.NewReservationService(reservationRepo, repo)
	reservationHandler := handler.NewReservationHandler(reservationService)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(corsmiddleware.Cors)

	// Отдача статических файлов (картинок)
	imgDir := "/usr/share/nginx/html/images"
	if _, err := os.Stat(imgDir); os.IsNotExist(err) {
		if err := os.MkdirAll(imgDir, 0755); err != nil {
			logger.Error("Failed to create images directory", zap.Error(err))
		}
	}
	fs := http.FileServer(http.Dir(imgDir))

	// Добавляем CORS для статических файлов
	corsFs := corsmiddleware.Cors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))

	r.Handle("/images/*", http.StripPrefix("/images/", corsFs))

	r.Route("/api/products", func(r chi.Router) {
		r.Post("/", productHandler.CreateProduct)
		r.Get("/", productHandler.ListProducts)
		r.Get("/{id}", productHandler.GetProduct)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
		r.Post("/{id}/availability", productHandler.CheckAvailability)
		r.Post("/{id}/image", productHandler.UploadImage)
	})

	// Регистрируем роуты для резервирования
	reservationHandler.Register(r)

	addr := cfg.Server.GetServerAddr()
	logger.Info("Starting server", zap.String("address", addr))
	if err := http.ListenAndServe(addr, r); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
