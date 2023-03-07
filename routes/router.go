package routes

import (
	"github.com/gin-gonic/gin"

	controllers "github.com/mrehanabbasi/appraisal-system-backend/controller"
)

func NewRouter() *gin.Engine {

	router := gin.Default()

	roleController := controllers.NewRoleController()

	v1 := router.Group("/v1")
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

	return router
}
