package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/mrehanabbasi/appraisal-system-backend/service"
)

func NewRouter() *gin.Engine {

	sc := service.NewSupervisorService()

	router := gin.Default()
	ec := service.NewEmployeeService()

	kc := service.NewKPIService()

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
	roleController := service.NewRoleService()

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

	v1 = router.Group("/v1")
	{
		roles := v1.Group("/kpis")
		{

			roles.POST("/", kc.CreateKPI)
			roles.GET("/", kc.GetAllKPI)
			roles.GET(":id", kc.GetKPIByID)
			roles.PUT(":id", kc.UpdateKPI)

		}
	}
	return router
}
