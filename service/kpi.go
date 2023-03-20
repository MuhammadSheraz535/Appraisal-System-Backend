package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mrehanabbasi/appraisal-system-backend/controller"
	"github.com/mrehanabbasi/appraisal-system-backend/database"
	"github.com/mrehanabbasi/appraisal-system-backend/models"
	"gorm.io/gorm"
)

type KPIController struct {
	Db *gorm.DB
}

func NewKPIController() *KPIController {
	db := database.DB
	db.AutoMigrate(&models.KPI{}, models.MeasuredData{}, models.QuestionaireData{})
	return &KPIController{Db: db}
}

func (r *KPIController) CreateKPI(c *gin.Context) {
	var kpi models.KPI
	err := c.ShouldBindJSON(&kpi)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	kpi, err = controller.CreateKPI(r.Db, kpi)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, kpi)
}

func (r *KPIController) GetAllKPI(c *gin.Context) {
	queryParams := make(map[string]string)

	if kpiName := c.Query("kpi_name"); kpiName != "" {
		queryParams["kpi_name"] = kpiName
	}

	if assignType := c.Query("assign_type"); assignType != "" {
		queryParams["assign_type"] = assignType
	}

	if rolesApplicable := c.Query("roles_applicable"); rolesApplicable != "" {
		queryParams["roles_applicable"] = rolesApplicable
	}

	kpis, err := controller.GetAllKPI(r.Db, queryParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, kpis)
}

func (r *KPIController) GetKPIByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	kpi, err := controller.GetKPIByID(r.Db, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, kpi)
}

func (r *KPIController) UpdateKPI(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var updatedKPI models.KPI
	err = c.ShouldBindJSON(&updatedKPI)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	kpi, err := controller.UpdateKPI(r.Db, uint(id), updatedKPI)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, kpi)
}

func (r *KPIController) DeleteKPI(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := controller.DeleteKPI(r.Db, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
