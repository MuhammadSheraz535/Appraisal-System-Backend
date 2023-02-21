package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mrehanabbasi/appraisal-system-backend/controller"
)

func NewRouter() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/v1")

	uc := controller.New()
	v1.POST("/employees", uc.CreateEmployee)
	v1.GET("/employees", uc.GetEmployees)
	v1.GET("/employees/:id", uc.GetEmployee)
	v1.GET("/employee/:name", uc.GetEmployeeByName)
	v1.PUT("/employees/:id", uc.UpdateEmployee)
	v1.DELETE("/employees/:id", uc.DeleteEmployee)
	v1.GET("/supervisor/:name", uc.GetSupervisorByName)
	v1.GET("/role/:name", uc.GetByRole)

	return router
}
