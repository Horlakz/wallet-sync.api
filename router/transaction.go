package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/horlakz/wallet-sync.api/handler"
	"github.com/horlakz/wallet-sync.api/internal/config"
	"github.com/horlakz/wallet-sync.api/lib/database"
	"github.com/horlakz/wallet-sync.api/middleware"
	core_repository "github.com/horlakz/wallet-sync.api/repository/core"
	"github.com/horlakz/wallet-sync.api/service"
)

func InitializeTransactionRouter(router fiber.Router, db database.DatabaseInterface, env config.Env) {
	// Repositories
	transactionRepository := core_repository.NewTransactionRepository(db)

	// Services
	transactionService := service.NewTransactionService(transactionRepository)

	// Handlers
	transactionHandler := handler.NewTransactionHandler(transactionService)

	// middlewares
	authMiddleware := middleware.Protected()

	// Base routes
	transactionRoute := router.Group("/transaction", authMiddleware)

	// Routes
	transactionRoute.Get("/", transactionHandler.GetTransactionHistory)
}
