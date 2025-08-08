package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/horlakz/wallet-sync.api/handler"
	"github.com/horlakz/wallet-sync.api/internal/config"
	"github.com/horlakz/wallet-sync.api/lib/database"
	user_repository "github.com/horlakz/wallet-sync.api/repository/user"
	"github.com/horlakz/wallet-sync.api/service"
)

func InitializeUserRouter(router fiber.Router, db database.DatabaseInterface, env config.Env) {
	// Repositories
	userRepository := user_repository.NewUserRepository(db)

	// Services
	authService := service.NewAuthService(userRepository)

	// Handler
	authHandler := handler.NewAuthHandler(authService)

	// Routers
	authRoute := router.Group("/auth")

	// Routes
	authRoute.Post("/login", authHandler.Login)
	authRoute.Post("/register", authHandler.Register)
}
