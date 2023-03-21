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

type KPIService struct {
	Db *gorm.DB
}

func NewKPIService() *KPIService {
	db := database.DB
	db.AutoMigrate(&models.KPI{}, models.MeasuredData{}, models.QuestionaireData{})
	return &KPIService{Db: db}
}

func (r *KPIService) CreateKPI(c *gin.Context) {
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

func (r *KPIService) GetAllKPI(c *gin.Context) {
	queryParams := make(map[string]string)

	if kpiName := c.Query("kpi_name"); kpiName != "" {
		queryParams["kpi_name"] = kpiName
	}

	if assignType := c.Query("assign_type"); assignType != "" {
		queryParams["assign_type"] = assignType
	}

	if ApplicableFor := c.Query("applicable_for"); ApplicableFor != "" {
		queryParams["applicable_for"] = ApplicableFor
	}

	kpis, err := controller.GetAllKPI(r.Db, queryParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, kpis)
}

func (r *KPIService) GetKPIByID(c *gin.Context) {
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

func (r *KPIService) UpdateKPI(c *gin.Context) {
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

func (r *KPIService) DeleteKPI(c *gin.Context) {
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
