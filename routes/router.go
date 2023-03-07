package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/mrehanabbasi/appraisal-system-backend/controller"

	controllers "github.com/mrehanabbasi/appraisal-system-backend/controller"
)

func NewRouter() *gin.Engine {

	sc := controller.NewSupervisorController()

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
	roleController := controller.NewRoleController()

	v1 = router.Group("/v1")
	{
		roles := v1.Group("/roles")
		{
			roles.GET("/", roleController.GetAllRoles)
			roles.GET(":id", roleController.GetRoleByID)
			roles.POST("/", roleController.CreateRole)
			roles.PUT(":id", roleController.UpdateRole)
			roles.DELETE(":id", roleController.DeleteRole)
		}
	}

	v1 = router.Group("/v1")
	{
		supervisors := v1.Group("/supervisors")
		{
			supervisors.POST("/", sc.ConvertSupervisorToEmployee)
			supervisors.GET("/", sc.GetSupervisors)
			supervisors.GET(":id", sc.GetSupervisorById)
			supervisors.PUT(":id", sc.UpdateSupervisor)
			supervisors.DELETE(":id", sc.DeleteSupervisor)
		}
	}
	return router
}
