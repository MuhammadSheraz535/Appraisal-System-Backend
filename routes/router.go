package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mrehanabbasi/appraisal-system-backend/controller"
)

func NewRouter() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/v1")
	// TODO: Define/Call routes here
	sc := controller.NewSupervisorController()

	v1.POST("/supervisors", sc.CreateSupervisor)
	v1.GET("/supervisors", sc.GetSupervisors)
	v1.GET("/supervisors/:id", sc.GetSupervisorByID)
	v1.PUT("/supervisors/:id", sc.UpdateSupervisor)
	v1.DELETE("/supervisors/:id", sc.DeleteEmployee)
	v1.GET("/supervisor/:name", sc.GetSupervisorByName)

	return router
}
