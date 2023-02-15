package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mrehanabbasi/appraisal-system-backend/controller"
)

func UserRoutes(router *gin.Engine) {
	router.GET("/", controller.UserController)
}
