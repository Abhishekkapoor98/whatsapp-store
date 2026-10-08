package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"whatsapp-store/internal/handler"
	"whatsapp-store/internal/repository"
	"whatsapp-store/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

func main() {
	// 1. Connect to Database (using port 6432 as resolved earlier)
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://admin:password@localhost:6432/store_db?sslmode=disable"
	}

	// Ensure sslmode is required for external connections (like Render)
	// If using an internal Render URL, this might need to be removed or handled differently.
	// However, usually, Render requires sslmode=require for connections.
	if !strings.Contains(dsn, "sslmode=") {
		if strings.Contains(dsn, "localhost") {
			dsn += "?sslmode=disable"
		} else {
			// Determine if it needs ? or &
			if strings.Contains(dsn, "?") {
				dsn += "&sslmode=require"
			} else {
				dsn += "?sslmode=require"
			}
		}
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}
	defer db.Close()

	// 2. Initialize sqlc Queries
	queries := repository.New(db)

	// 3. Initialize Repositories
	adminRepo := repository.NewAdminRepository(queries)
	productRepo := repository.NewProductRepository(queries)

	// 4. Initialize Usecases
	adminUsecase := usecase.NewAdminUsecase(adminRepo, "my-super-secret-key")
	productUsecase := usecase.NewProductUsecase(productRepo)

	// 5. Initialize Handlers
	adminHandler := handler.NewAdminHandler(adminUsecase)
	productHandler := handler.NewProductHandler(productUsecase)

	// 6. Setup Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Serve Static Files (For uploaded images)
	fs := http.FileServer(http.Dir("web/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Serve HTML Pages
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("web/templates/index.html")
		if err != nil {
			http.Error(w, "Could not load template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	r.Get("/admin/dashboard", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("web/templates/admin.html")
		if err != nil {
			http.Error(w, "Could not load template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	// Public API Routes
	r.Route("/api", func(r chi.Router) {
		r.Get("/products", productHandler.GetAvailableProducts)
	})

	// Admin API Routes
	r.Route("/admin", func(r chi.Router) {
		r.Post("/login", adminHandler.Login)
		r.Get("/products", productHandler.GetAllProducts) // Added for Inventory Management
		r.Post("/products", productHandler.AddProduct)
		r.Put("/products/{id}/stock", productHandler.UpdateStock)
	})

	// 7. Start Server
	log.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
