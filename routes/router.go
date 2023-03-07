package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/mrehanabbasi/appraisal-system-backend/controller"
)

func NewRouter() *gin.Engine {
	router := gin.Default()
	ec := controllers.NewEmployeeController()

	v1 := router.Group("/v1")
	{
		employee := v1.Group("/employees")
		{
			employee.POST("", ec.CreateEmployee)
			employee.GET("", ec.GetEmployees)
			employee.GET("/:id", ec.GetEmployee)
			employee.PUT("/:id", ec.UpdateEmployee)
			employee.DELETE("/:id", ec.DeleteEmployee)
		}
	}
	return router
}
