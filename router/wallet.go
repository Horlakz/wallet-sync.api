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

func InitializeWalletRouter(router fiber.Router, db database.DatabaseInterface, env config.Env) {
	// Repositories
	accountRepository := core_repository.NewAccountRepository(db)
	transactionRepository := core_repository.NewTransactionRepository(db)
	ledgerEntryRepository := core_repository.NewLedgerEntryRepository(db)

	// Services
	walletService := service.NewWalletService(accountRepository, transactionRepository, ledgerEntryRepository, db)

	// Handlers
	walletHandler := handler.NewWalletHandler(walletService)

	// middlewares
	authMiddleware := middleware.Protected()

	// Base routes
	walletRoute := router.Group("/wallet", authMiddleware)

	// Routes
	walletRoute.Get("/", walletHandler.GetDetails)
	walletRoute.Post("/fund", walletHandler.Fund)
	walletRoute.Post("/withdraw", walletHandler.Withdraw)
	walletRoute.Post("/transfer", walletHandler.Transfer)
}
