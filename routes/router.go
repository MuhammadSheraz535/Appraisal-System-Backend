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
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
			"Accept",
			"Origin",
			"Access-Control-Allow-Origin",
			"Cache-Control",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	ec := service.NewEmployeeService()
	roleController := service.NewRoleService()
	sc := service.NewSupervisorService()
	kc := service.NewKPIService()
	af := service.NewApprasialFlowService()

	v1 := router.Group("/v1")

	employee := v1.Group("/employees")
	{
		employee.POST("", ec.CreateEmployee)
		employee.GET("", ec.GetEmployees)
		employee.GET("/:id", ec.GetEmployee)
		employee.PUT("/:id", ec.UpdateEmployee)
		employee.DELETE("/:id", ec.DeleteEmployee)
	}

	roles := v1.Group("/roles")
	{
		roles.GET("/", roleController.GetAllRoles)
		roles.GET(":id", roleController.GetRoleByID)
		roles.POST("/", roleController.CreateRole)
		roles.PUT(":id", roleController.UpdateRole)
		roles.DELETE(":id", roleController.DeleteRole)
	}

	supervisors := v1.Group("/supervisors")
	{
		supervisors.POST("", sc.ConvertSupervisorToEmployee)
		supervisors.GET("", sc.GetSupervisors)
		supervisors.GET("/:id", sc.GetSupervisorById)
		supervisors.PUT("/:id", sc.UpdateSupervisor)
		supervisors.DELETE("/:id", sc.DeleteSupervisor)
	}

	kpis := v1.Group("/kpis")
	{
		kpis.POST("", kc.CreateKPI)
		kpis.GET("", kc.GetAllKPI)
		kpis.GET("/:id", kc.GetKPIByID)
		kpis.PUT("/:id", kc.UpdateKPI)
		kpis.DELETE("/:id", kc.DeleteKPI)
	}

	appraisalFlows := v1.Group("/appraisal_flows")
	{
		appraisalFlows.POST("", af.CreateAppraisalFlow)
		appraisalFlows.GET("", af.GetAllApprasialFlow)
		appraisalFlows.GET("/:id", af.GetAppraisalFlowByID)
		appraisalFlows.PUT("/:id", af.UpdateAppraisalFlow)
		appraisalFlows.DELETE("/:id", af.DeleteApprasialFlow)
	}

	return router
}
