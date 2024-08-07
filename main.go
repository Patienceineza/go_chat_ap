package main

import (
	"chat-app-backend/internal/config"
	"chat-app-backend/internal/database"
	
	
	"chat-app-backend/internal/router"
	seed "chat-app-backend/internal/seeds"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.LoadConfig()
	db, err := database.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	seed.Seed(db)

	r := router.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
