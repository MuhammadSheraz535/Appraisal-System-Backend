package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mrehanabbasi/appraisal-system-backend/database"
)

func main() {
	router := gin.Default()
	database.Connect()

	router.Run(":8080")
}
