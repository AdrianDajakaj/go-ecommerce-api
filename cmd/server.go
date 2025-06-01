package main

import (
	"go-ecommerce-api/internal/infrastructure/persistence/sqlite"
	"go-ecommerce-api/internal/interface/http"
	"log"

	"github.com/joho/godotenv"
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
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
