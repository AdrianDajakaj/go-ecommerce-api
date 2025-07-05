package main

import (
	"go-ecommerce-api/internal/infrastructure/persistence/sqlite"
	httpRouter "go-ecommerce-api/internal/interface/http"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Missing .env file or error while loading")
	}

	// Determine database path - support both local and Docker environments
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "ecommerce.db"
		// Check if we're running in Docker container
		if _, err := os.Stat("/app/data"); err == nil {
			dbPath = "/app/data/ecommerce.db"
		}
	}

	// Initialize database
	db, err := sqlite.NewGormDB(dbPath)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	// Create Echo router
	e := httpRouter.NewRouter(db)

	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	// Health check endpoint for Docker
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "healthy",
			"service": "go-ecommerce-api",
		})
	})

	// Static assets serving - support both local and Docker environments
	assetsPath := os.Getenv("ASSETS_PATH")
	if assetsPath == "" {
		assetsPath = "assets"
		// Check if we're running in Docker container
		if _, err := os.Stat("/app/assets"); err == nil {
			assetsPath = "/app/assets"
		}
	}

	// Serve static files
	if _, err := os.Stat(assetsPath); err == nil {
		e.Static("/assets", assetsPath)
		log.Printf("Serving static assets from: %s", assetsPath)
	} else {
		log.Printf("Assets directory not found: %s", assetsPath)
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	log.Printf("Database path: %s", dbPath)
	log.Printf("Assets path: %s", assetsPath)

	if err := e.Start(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
