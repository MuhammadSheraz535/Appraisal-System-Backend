package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mrehanabbasi/appraisal-system-backend/controller"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/v1")

	// TODO:Call routes here

	uc := controller.New()
	v1.POST("/employees", uc.CreateUser)
	v1.GET("/employees", uc.GetUsers)
	v1.GET("/employees/:id", uc.GetUser)
	v1.GET("/employee/:name", uc.GetUserByName)
	v1.PUT("/employees/:id", uc.UpdateUser)
	v1.DELETE("/employees/:id", uc.DeleteUser)
	v1.GET("/supervisor/:name", uc.GetSupervisorByName)
	v1.GET("/role/:name", uc.GetByRole)

	return r
}
