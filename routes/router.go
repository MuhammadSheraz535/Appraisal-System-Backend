package routes

import "github.com/gin-gonic/gin"

func Allroutes() {
	router := gin.Default()
	router.Run(":8080")
}
