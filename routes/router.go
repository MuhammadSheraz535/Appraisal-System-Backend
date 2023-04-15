package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/cors"
	"github.com/mrehanabbasi/appraisal-system-backend/service"
)

func NewRouter() *gin.Engine {
	router := gin.Default()

	_ = router.SetTrustedProxies(nil)

	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "Accept", "Origin", "Cache-Control"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	sc := service.NewSupervisorService()

	ec := service.NewEmployeeService()

	kc := service.NewKPIService()

	af := service.NewApprasialFlowService()

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
			roles.DELETE(":id", kc.DeleteKPI)

		}
	}

	v1 = router.Group("/v1")
	{
		appraisalflow := v1.Group("/appraisal_flows")
		{

			appraisalflow.POST("/", af.CreateAppraisalFlow)
			appraisalflow.GET("/", af.GetAllApprasialFlow)
			appraisalflow.GET(":id", af.GetAppraisalFlowByID)
			appraisalflow.PUT(":id", af.UpdateAppraisalFlow)
			// appraisalflow.DELETE(":id", kc.DeleteKPI)

		}
	}
	return router
}
