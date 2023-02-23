package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mrehanabbasi/appraisal-system-backend/controller"
)

func NewRouter() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/v1/employees")

	ec := controller.NewEmployeeController()
	v1.POST("", ec.CreateEmployee)
	v1.GET("", ec.GetEmployees)
	v1.GET("/:id", ec.GetEmployee)
	v1.PUT("/:id", ec.UpdateEmployee)
	v1.DELETE("/:id", ec.DeleteEmployee)

	return router
}
