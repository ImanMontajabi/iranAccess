package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"iranAccess/internal/checker"
	"iranAccess/internal/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

const (
	defaultPort    = ":3000"
	csvFilePath    = "./domains.csv"
	staticFilePath = "./public"
)

func main() {
	// Create domain checker
	domainChecker := checker.NewDomainChecker(csvFilePath)

	// Start domain checking in background
	go domainChecker.StartDomainChecker()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Iran Access Checker",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// Routes
	setupRoutes(app)

	// Graceful shutdown
	go func() {
		fmt.Printf("Server is running on port %s\n", defaultPort)
		if err := app.Listen(defaultPort); err != nil {
			log.Fatal("Server failed to start:", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	fmt.Println("Server exited")
}

func setupRoutes(app *fiber.App) {
	// API routes
	api := app.Group("/api")
	api.Get("/check", handlers.CheckDomainsHandler)

	// Legacy route for backward compatibility
	app.Get("/check", handlers.CheckDomainsHandler)

	// Static files
	app.Static("/", staticFilePath)
}
