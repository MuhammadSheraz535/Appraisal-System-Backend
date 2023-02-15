package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mrehanabbasi/appraisal-system-backend/database"
	"github.com/mrehanabbasi/appraisal-system-backend/routes"
)

func main() {
	router := gin.Default()
	database.Connect()
	routes.UserRoutes(router)

	router.Run(":8080")
}
