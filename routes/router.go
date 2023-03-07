package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/mrehanabbasi/appraisal-system-backend/controller"
)

func NewRouter() *gin.Engine {
	router := gin.Default()

	// TODO: Define/Call routes here
	sc := controllers.NewSupervisorController()

	v1 := router.Group("/v1")
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
