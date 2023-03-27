package service

import (
	"errors"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/mrehanabbasi/appraisal-system-backend/database"
	"github.com/mrehanabbasi/appraisal-system-backend/models"
	"gorm.io/gorm"
)

type KPIService struct {
	Db *gorm.DB
}

func NewKPIService() *KPIService {
	db := database.DB
	db.AutoMigrate(&models.KpiType{}, &models.FeedbackKpi{}, &models.MeasuredKpi{}, &models.ObservatoryKpi{}, &models.QuestionaireKpi{})
	return &KPIService{Db: db}
}

const (
	FEEDBACK_KPI_TYPE      = "Feedback"
	OBSERVATORY_KPI_TYPE   = "Observatory"
	MEASURED_KPI_TYPE      = "Measured"
	QUESTIONNAIRE_KPI_TYPE = "Questionnaire"
)

func (s *KPIService) CreateKPI(c *gin.Context) {
	var kpiCommon models.KpiCommon
	err := c.ShouldBindBodyWith(&kpiCommon, binding.JSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid KPI data"})
		return
	}
	switch kpiCommon.KpiType {
	case FEEDBACK_KPI_TYPE:
		var feedbackKpi models.FeedbackKpi
		err := c.ShouldBindBodyWith(&feedbackKpi, binding.JSON)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid feedback KPI data"})
			return
		}

		feedbackKpi.KpiCommon = kpiCommon
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

		observatoryKpi.KpiCommon = kpiCommon
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

		measuredKpi.KpiCommon = kpiCommon
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

		questionnaireKpi.KpiCommon = kpiCommon
		result := s.Db.Create(&questionnaireKpi)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, questionnaireKpi)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid KPI type"})
		return
	}
}

func (s *KPIService) GetAllKPI(c *gin.Context) {
	// Get query parameters
	kpiName := c.Query("KpiName")
	assignType := c.Query("AssignType")
	kpiType := c.Query("KpiType")

	// Check if kpiType query parameter is present
	if kpiType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "KpiType query parameter is required"})
		return
	}

	var kpis interface{}

	// Get KPI records based on kpiType query parameter
	switch kpiType {
	case "feedback":
		var feedbackKpis []models.FeedbackKpi
		db := s.Db.Where("kpi_type = ?", kpiType)
		if kpiName != "" {
			db = db.Where("kpi_name = ?", kpiName)
		}
		if assignType != "" {
			db = db.Where("assign_type = ?", assignType)
		}
		if err := db.Find(&feedbackKpis).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Feedback KPIs"})
			return
		}
		kpis = feedbackKpis
	case "observatory":
		var observatoryKpis []models.ObservatoryKpi
		db := s.Db.Where("kpi_type = ?", kpiType)
		if kpiName != "" {
			db = db.Where("kpi_name = ?", kpiName)
		}
		if assignType != "" {
			db = db.Where("assign_type = ?", assignType)
		}
		if err := db.Find(&observatoryKpis).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Observatory KPIs"})
			return
		}
		kpis = observatoryKpis
	case "questionnaire":
		var questionnaireKpis []models.QuestionaireKpi
		db := s.Db.Where("kpi_type = ?", kpiType)
		if kpiName != "" {
			db = db.Where("kpi_name = ?", kpiName)
		}
		if assignType != "" {
			db = db.Where("assign_type = ?", assignType)
		}
		if err := db.Find(&questionnaireKpis).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Questionnaire KPIs"})
			return
		}
		kpis = questionnaireKpis
	case "measured":
		var measuredKpis []models.MeasuredKpi
		db := s.Db.Where("kpi_type = ?", kpiType)
		if kpiName != "" {
			db = db.Where("kpi_name = ?", kpiName)
		}
		if assignType != "" {
			db = db.Where("assign_type = ?", assignType)
		}
		if err := db.Find(&measuredKpis).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Measured KPIs"})
			return
		}
		kpis = measuredKpis
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid KpiType query parameter"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": kpis})
}

func (s *KPIService) GetKPIByID(c *gin.Context) {
	kpiType := c.Query("KpiType")
	kpiID := c.Param("id")

	if kpiType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "KpiType query parameter is required"})
		return
	}

	var kpi interface{}

	switch kpiType {
	case "feedback":
		var feedbackKpi models.FeedbackKpi
		if err := s.Db.Where("id = ?", kpiID).First(&feedbackKpi).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Feedback KPI not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Feedback KPI"})
			return
		}
		kpi = feedbackKpi
	case "observatory":
		var observatoryKpi models.ObservatoryKpi
		if err := s.Db.Where("id = ?", kpiID).First(&observatoryKpi).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Observatory KPI not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Observatory KPI"})
			return
		}
		kpi = observatoryKpi
	case "questionnaire":
		var questionnaireKpi models.QuestionaireKpi
		if err := s.Db.Where("id = ?", kpiID).First(&questionnaireKpi).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Questionnaire KPI not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Questionnaire KPI"})
			return
		}
		kpi = questionnaireKpi
	case "measured":
		var measuredKpi models.MeasuredKpi
		if err := s.Db.Where("id = ?", kpiID).First(&measuredKpi).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Measured KPI not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Measured KPI"})
			return
		}
		kpi = measuredKpi
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid KpiType query parameter"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": kpi})
}

func (s *KPIService) UpdateKPI(c *gin.Context) {
	kpiType := c.Query("KpiType")
	kpiID := c.Param("id")

	if kpiType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "KpiType query parameter is required"})
		return
	}

	switch kpiType {
	case "feedback":
		var feedbackKpi models.FeedbackKpi
		if err := s.Db.Where("id = ?", kpiID).First(&feedbackKpi).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Feedback KPI not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Feedback KPI"})
			return
		}

		// Bind the request body to the FeedbackKpi model and validate it
		if err := c.BindJSON(&feedbackKpi); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update the FeedbackKpi model in the database
		if err := s.Db.Save(&feedbackKpi).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Feedback KPI"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": feedbackKpi})

	case "observatory":
		var observatoryKpi models.ObservatoryKpi
		if err := s.Db.Where("id = ?", kpiID).First(&observatoryKpi).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Observatory KPI not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Observatory KPI"})
			return
		}

		// Bind the request body to the ObservatoryKpi model and validate it
		if err := c.BindJSON(&observatoryKpi); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Update the ObservatoryKpi model in the database
		if err := s.Db.Save(&observatoryKpi).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Observatory KPI"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": observatoryKpi})

	case "questionnaire":
		var questionnaireKpi models.QuestionaireKpi
		if err := s.Db.Where("id = ?", kpiID).First(&questionnaireKpi).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Questionnaire KPI not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Questionnaire KPI"})
			return
		}

		// Bind the request body to the QuestionnaireKpi model and validate it
		if err := c.BindJSON(&questionnaireKpi); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update the QuestionnaireKpi model in the database
		if err := s.Db.Save(&questionnaireKpi).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Questionnaire KPI"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": questionnaireKpi})

	case "measured":
		var measuredKpi models.MeasuredKpi
		if err := s.Db.Where("id = ?", kpiID).First(&measuredKpi).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Measured KPI not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Measured KPI"})
			return
		}

		// Bind the request body to the MeasuredKpi model and validate it
		if err := c.BindJSON(&measuredKpi); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update the MeasuredKpi model in the database
		if err := s.Db.Save(&measuredKpi).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Measured KPI"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": measuredKpi})

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid KpiType query parameter"})
		return
	}

}

func (s *KPIService) DeleteKPI(c *gin.Context) {
	kpiType := c.Query("KpiType")
	kpiID := c.Param("id")

	if kpiType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "KpiType query parameter is required"})
		return
	}

	switch kpiType {
	case "feedback":
		var feedbackKpi models.FeedbackKpi
		if err := s.Db.Where("id = ?", kpiID).First(&feedbackKpi).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Feedback KPI not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Feedback KPI"})
			return
		}

		if err := s.Db.Delete(&feedbackKpi).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Feedback KPI"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Feedback KPI deleted successfully"})

	case "observatory":
		var observatoryKpi models.ObservatoryKpi
		if err := s.Db.Where("id = ?", kpiID).First(&observatoryKpi).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Observatory KPI not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Observatory KPI"})
			return
		}

		if err := s.Db.Delete(&observatoryKpi).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Observatory KPI"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Observatory KPI deleted successfully"})

	case "questionnaire":
		var questionnaireKpi models.QuestionaireKpi
		if err := s.Db.Where("id = ?", kpiID).First(&questionnaireKpi).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Questionnaire KPI not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Questionnaire KPI"})
			return
		}

		if err := s.Db.Delete(&questionnaireKpi).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Questionnaire KPI"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Questionnaire KPI deleted successfully"})

	case "measured":
		var measuredKpi models.MeasuredKpi
		if err := s.Db.Where("id = ?", kpiID).First(&measuredKpi).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Measured KPI not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Measured KPI"})
			return
		}

		if err := s.Db.Delete(&measuredKpi).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Measured KPI"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Measured KPI deleted successfully"})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid KPI type"})
	}
}
