package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mrehanabbasi/appraisal-system-backend/service"
)

func NewRouter() *gin.Engine {

	router := gin.Default()

	KPIController := service.NewKPIController()

	v1 := router.Group("/v1")
	{
		roles := v1.Group("/kpis")
		{
			roles.GET("/", KPIController.GetAllKPI)
			roles.GET(":id", KPIController.GetKPIByID)
			roles.POST("/", KPIController.CreateKPI)
			roles.PUT(":id", KPIController.UpdateKPI)
			roles.DELETE(":id", KPIController.DeleteKPI)
		}
	}

	return router
}
