package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/Cakra17/imphnen/docs"
	"github.com/Cakra17/imphnen/internal/config"
	"github.com/Cakra17/imphnen/internal/handlers"
	md "github.com/Cakra17/imphnen/internal/middleware"
	"github.com/Cakra17/imphnen/internal/store"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Imphnen API
// @version         0.1
// @description     API for managing users, products, and orders
// @termsOfService  http://swagger.io/terms/

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @tag.name Auth
// @tag.description Authentication endpoints for user registration and login

// @tag.name Users
// @tag.description Operations related to user management
// @tag.docs.url https://example.com/docs/users
// @tag.docs.description User management documentation

// @tag.name Receipts
// @tag.description Operations related to receipt
// @tag.docs.url https://example.com/docs/receipts

// @tag.name Transactions
// @tag.description Operations related to order processing

func main() {
	cfg := config.Load()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(middleware.Timeout(time.Minute))

	server := &http.Server{
		Addr:         cfg.Port,
		Handler:      r,
		WriteTimeout: time.Minute,
		ReadTimeout:  time.Minute,
		IdleTimeout:  time.Minute,
	}

	db := config.ConnectDB(cfg.DSN)
	cld, err := cloudinary.NewFromParams(cfg.CloudinaryName, cfg.CloudinaryApiKey, cfg.CLoudinaryApiSecret)
	if err != nil {
		log.Fatalf("Failed connect to cloudinary service: %v", err)
	}

	userRepo := store.NewUserRepo(db)
	receiptRepo := store.NewReceiptRepo(db)
	transactionRepo := store.NewTransactionRepo(db)

	userHandler := handlers.NewUserHandler(handlers.UserHandlerConfig{
		UserRepo:      userRepo,
		JwtSecret:     cfg.JWTSecret,
		TokenDuration: time.Hour * 8,
	})

	receiptHandler := handlers.NewReceiptHandler(handlers.ReceiptHandlerConfig{
		ReceiptRepo:     receiptRepo,
		TransactionRepo: transactionRepo,
		Cld:             cld,
	})

	transactionHandler := handlers.NewTransactionHandler(handlers.TransactionHandlerConfig{
		TransactionStore: &transactionRepo,
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/docs/*", httpSwagger.WrapHandler)

		r.Post("/auth/login", userHandler.Login)
		r.Post("/auth/register", userHandler.Register)

		r.Route("/users", func(r chi.Router) {
			r.Use(md.Auth)
			r.Get("/me", userHandler.Session)
		})

		r.Route("/receipts", func(r chi.Router) {
			r.Use(md.Auth)
			r.Post("/", receiptHandler.CreateReceipt)
			r.Get("/", receiptHandler.GetReceipts)
			r.Get("/{id}", receiptHandler.GetReceiptByID)
			r.Get("/items/{id}", receiptHandler.GetItemsByRecieptID)
		})

		r.Route("/transactions", func(r chi.Router) {
			r.Use(md.Auth)
			r.Post("/", transactionHandler.CreateTransaction)
			r.Get("/date", transactionHandler.GetTransactionsByDate)
			r.Get("/range", transactionHandler.GetTransactionsByRange)
			r.Get("/days", transactionHandler.GetTransactionsByDays)
			r.Get("/stats", transactionHandler.GetTransactionStats)
			r.Get("/stats/days", transactionHandler.GetTransactionStatsByDays)
			r.Get("/type", transactionHandler.GetTransactionsByType)
			r.Get("/source", transactionHandler.GetTransactionsBySource)
		})
	})

	closed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		log.Println("Received shutdown signal, shutting down server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Failed to shutdown server: %v", err)
		}

		close(closed)
	}()

	log.Printf("server running on port %s", server.Addr[1:])
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Failed to run server: %v", err)
	}

	<-closed
	log.Println("Server shutdown gracefully")
}
