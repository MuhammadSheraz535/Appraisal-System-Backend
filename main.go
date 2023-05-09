package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/mrehanabbasi/appraisal-system-backend/database"
	"github.com/mrehanabbasi/appraisal-system-backend/logger"
	"github.com/mrehanabbasi/appraisal-system-backend/routes"
)

func main() {
	// Load environment variables from .env file
	_ = godotenv.Load(".env")

	// Initializing logger
	logger.TextLogInit()

	// Connect to the database
	database.Connect()

	// Register all the routes
	server := routes.NewRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	_ = server.Run(":" + port)
}
