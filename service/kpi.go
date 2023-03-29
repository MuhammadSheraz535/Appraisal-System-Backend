// Service/kpi.go

package service

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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
	db.AutoMigrate(&models.Kpis{}, &models.KpiType{}, &models.FeedbackKpi{}, &models.MeasuredKpi{}, &models.ObservatoryKpi{}, &models.QuestionaireKpi{})
	return &KPIService{Db: db}
}

const (
	FEEDBACK_KPI_TYPE      = "Feedback"
	OBSERVATORY_KPI_TYPE   = "Observatory"
	MEASURED_KPI_TYPE      = "Measured"
	QUESTIONNAIRE_KPI_TYPE = "Questionnaire"
)

func (s *KPIService) CreateKPI(c *gin.Context) {
	var kpi models.Kpis

	if err := c.ShouldBindBodyWith(&kpi, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	kpi.ID = 0
	switch kpi.KpiType {
	case FEEDBACK_KPI_TYPE:
		var feedbackKpi models.FeedbackKpi
		err := c.ShouldBindBodyWith(&feedbackKpi, binding.JSON)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid feedback KPI data"})
			return
		}

		result := s.Db.Create(&feedbackKpi)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, feedbackKpi)
	case OBSERVATORY_KPI_TYPE:
		var observatoryKpi models.ObservatoryKpi
		err := c.ShouldBindBodyWith(&observatoryKpi, binding.JSON)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid observatory KPI data"})
			return
		}
		result := s.Db.Create(&observatoryKpi)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, observatoryKpi)
	case MEASURED_KPI_TYPE:
		var measuredKpi models.MeasuredKpi
		err := c.ShouldBindBodyWith(&measuredKpi, binding.JSON)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid measured KPI data"})
			return
		}
		result := s.Db.Create(&measuredKpi)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, measuredKpi)
	case QUESTIONNAIRE_KPI_TYPE:
		var questionnaireKpi models.QuestionaireKpi
		err := c.ShouldBindBodyWith(&questionnaireKpi, binding.JSON)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid questionnaire KPI data"})
			return
		}
		result := s.Db.Create(&questionnaireKpi)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, questionnaireKpi)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type"})
		return
	}

	kpi, err := controller.CreateKPI(s.Db, kpi)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, kpi)
}

func (s *KPIService) GetAllKPI(c *gin.Context) {
	var kpis []models.Kpis
	db := s.Db

	kpiName := c.Query("KpiName")
	assignType := c.Query("AssignType")
	kpiType := c.Query("KpiType")

	if kpiName != "" && assignType != "" && kpiType != "" {
		db = db.Where("kpis.kpi_name = ? AND kpis.assign_type = ? AND kpis.kpi_type = ?", kpiName, assignType, kpiType)
	} else if kpiName != "" && assignType != "" {
		db = db.Where("kpis.kpi_name = ? AND kpis.assign_type = ?", kpiName, assignType)
	} else if kpiName != "" && kpiType != "" {
		db = db.Where("kpis.kpi_name = ? AND kpis.kpi_type = ?", kpiName, kpiType)
	} else if assignType != "" && kpiType != "" {
		db = db.Where("kpis.assign_type = ? AND kpis.kpi_type = ?", assignType, kpiType)
	} else if kpiName != "" {
		db = db.Where("kpis.kpi_name = ?", kpiName)
	} else if assignType != "" {
		db = db.Where("kpis.assign_type = ?", assignType)
	} else if kpiType != "" {
		db = db.Where("kpis.kpi_type = ?", kpiType)
	}

	if err := controller.GetAllKPI(db, &kpis); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch KPIs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": kpis})
}

func (s *KPIService) GetKPIByID(c *gin.Context) {
	kpiID := c.Param("id")
	var kpi models.Kpis

	if err := controller.GetKPIByID(s.Db, kpiID, &kpi); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "KPI not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch KPI"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": kpi})
}

func (s *KPIService) UpdateKPI(c *gin.Context) {
	kpiID := c.Param("id")
	var kpi models.Kpis

	if err := s.Db.Where("id = ?", kpiID).First(&kpi).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "KPI not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch KPI"})
		}
		return
	}

	if err := c.BindJSON(&kpi); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	kpi.ID = 0
	
	// Update the KPI
	kpi, err := controller.UpdateKPI(s.Db, kpi)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, kpi)
}

func (s *KPIService) DeleteKPI(c *gin.Context) {
	kpiID := c.Param("id")

	if err := controller.DeleteKPIByID(s.Db, kpiID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "KPI not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete KPI"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "KPI deleted successfully"})
}