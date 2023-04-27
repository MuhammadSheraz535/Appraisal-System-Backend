package routes

import (
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/cors"
	"github.com/mrehanabbasi/appraisal-system-backend/service"
)

func NewRouter() *gin.Engine {
	router := gin.Default()

	isCorsEnabled, _ := strconv.ParseBool(os.Getenv("ENABLE_CORS"))
	if isCorsEnabled {
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
				"Cache-Control",
			},
			ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
			AllowCredentials: false,
		}))
	}

	ec := service.NewEmployeeService()
	roleController := service.NewRoleService()
	sc := service.NewSupervisorService()
	kc := service.NewKPIService()
	af := service.NewApprasialFlowService()
	a := service.NewApprasialService()

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

	appraisal_flows := v1.Group("/appraisal_flows")
	{
		appraisal_flows.POST("", af.CreateAppraisalFlow)
		appraisal_flows.GET("", af.GetAllApprasialFlow)
		appraisal_flows.GET("/:id", af.GetAppraisalFlowByID)
		appraisal_flows.PUT("/:id", af.UpdateAppraisalFlow)
		appraisal_flows.DELETE("/:id", af.DeleteApprasialFlow)
	}

	appraisal := v1.Group("/appraisal")
	{
		appraisal.POST("", a.CreateAppraisal)
		appraisal.GET("", a.GetAllApprasial)
		appraisal.GET("/:id", a.GetAppraisalByID)
		appraisal.PUT("/:id", a.UpdateAppraisal)
		appraisal.DELETE("/:id", a.DeleteApprasial)
	}

	return router
}
