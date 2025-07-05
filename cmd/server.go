package main

import (
	"go-ecommerce-api/internal/infrastructure/persistence/sqlite"
	"go-ecommerce-api/internal/interface/http"
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	dsn := "ecommerce.db"
	db, err := sqlite.NewGormDB(dsn)
	if err := godotenv.Load(); err != nil {
		log.Println("Missing .env file or error while loading")
	}

	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	e := http.NewRouter(db)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	if err := e.Start(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
