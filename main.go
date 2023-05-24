package main

import (
	"os"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/mrehanabbasi/appraisal-system-backend/database"
	"github.com/mrehanabbasi/appraisal-system-backend/logger"
	"github.com/mrehanabbasi/appraisal-system-backend/routes"
)

func main() {
	// Load environment variables from .env file
	_ = godotenv.Load(".env")

	// Convert fe.Field() from StructField to json field for custom validation messages
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}

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
