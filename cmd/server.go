package main

import (
	"go-ecommerce-api/internal/infrastructure/persistence/sqlite"
	"go-ecommerce-api/internal/interface/http"
	"log"
)

func main() {
	dsn := "ecommerce.db"
	db, err := sqlite.NewGormDB(dsn)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	e := http.NewRouter(db)
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
