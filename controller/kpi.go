package controller

import (
	"net/http"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/mrehanabbasi/appraisal-system-backend/database"
	"github.com/mrehanabbasi/appraisal-system-backend/models"
)

type KPIController struct {
	Db *gorm.DB
}

func NewKPIController() *KPIController {
	db := database.DB
	db.AutoMigrate(&models.Role{})
	return &KPIController{Db: db}
}

func (r *KPIController) GetKPIs(c *gin.Context) {
	name := c.Query("name")
	assignType := c.Query("assign_type")
	kpiType := c.Query("type")
	role := c.Query("role")

	// Get KPIs based on query parameters
	kpis := []models.KPI{}
	query := r.Db
	if name != "" {
		query = query.Where("kpi_name = ?", name)
	}
	if assignType != "" {
		query = query.Where("assign_type = ?", assignType)
	}
	if kpiType != "" {
		query = query.Where("type = ?", kpiType)
	}
	if role != "" {
		query = query.Where("roles_applicable LIKE ?", "%"+role+"%")
	}
	if err := query.Find(&kpis).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, kpis)
}
