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

func (r *KPIController) GetKPIByID(c *gin.Context) {
	id := c.Param("id")

	// Get KPI by ID
	kpi := models.KPI{}
	if err := r.Db.Where("id = ?", id).First(&kpi).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, kpi)
}

// CreateKPI creates a new KPI
func (r *KPIController) CreateKPI(c *gin.Context) {
	kpi := models.KPI{}

	// Bind request body to KPI model
	if err := c.BindJSON(&kpi); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Create KPI
	if err := r.Db.Create(&kpi).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, kpi)
}

// UpdateKPI updates an existing KPI
func (r *KPIController) UpdateKPI(c *gin.Context) {
	id := c.Param("id")

	kpi := models.KPI{}

	// Get KPI by ID
	if err := r.Db.Where("id = ?", id).First(&kpi).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Bind request body to KPI model and update fields
	if err := c.BindJSON(&kpi); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Update KPI
	if err := r.Db.Save(&kpi).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, kpi)
}

func (r *KPIController) DeleteKPI(c *gin.Context) {
	id := c.Param("id")
	kpi := models.KPI{}
	// Get KPI by ID
	if err := r.Db.Where("id = ?", id).First(&kpi).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	// Delete KPI
	if err := r.Db.Delete(&kpi).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusNoContent)
}
