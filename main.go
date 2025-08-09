package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/horlakz/wallet-sync.api/internal/config"
	"github.com/horlakz/wallet-sync.api/internal/job"
	"github.com/horlakz/wallet-sync.api/internal/seed"
	"github.com/horlakz/wallet-sync.api/lib/database"
	"github.com/horlakz/wallet-sync.api/router"
)

func main() {
	app := fiber.New(fiber.Config{AppName: "Wallet Sync API v0.0.1"})

	app.Use(logger.New(logger.Config{}))
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
	app.Use(limiter.New(limiter.Config{
		Max:               1000,
		Expiration:        60 * time.Second,
		LimiterMiddleware: limiter.FixedWindow{},
	}))

	// Get environment variables
	env := config.GetEnv()

	// Start database connection
	dbConn := database.StartDatabaseClient(env)

	// Initialize router
	router.InitializeRouter(app, dbConn, env)

	// Migrate database
	database.Migrate(dbConn)

	// Seed database
	seed.NewSeeder(dbConn).Seed()

	// Initialize cron jobs
	cronService := job.NewCronService(dbConn)
	cronService.Start()

	log.Fatal(app.Listen("0.0.0.0:" + env.PORT))
}
