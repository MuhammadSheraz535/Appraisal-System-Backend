package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mrehanabbasi/appraisal-system-backend/controller"
)

func NewRouter() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/v1")

	// TODO: Define/Call routes here
	sc := controller.New()
	v1.POST("/supervisors", sc.CreateUser)
	v1.GET("/supervisors", sc.GetUsers)
	v1.GET("/supervisors/:id", sc.GetUser)
	v1.PUT("/supervisors/:id", sc.UpdateUser)
	v1.DELETE("/supervisors/:id", sc.DeleteUser)

	return router
}
