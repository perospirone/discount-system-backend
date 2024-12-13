package main

import (
	"discount-system-backend/internal/database"
	"discount-system-backend/internal/routes"
	"flag"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
}

func main() {
	migrate := flag.Bool("migrate", false, "Run database migrations")
	flag.Parse()

	db := database.Connection()

	if *migrate {
		log.Println("Running database migrations...")
		database.Migrate(db)
	}

	mux := routes.Routes()

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
