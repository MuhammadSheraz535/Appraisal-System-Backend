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

const (
	FEEDBACK_KPI_TYPE      = "Feedback"
	OBSERVATORY_KPI_TYPE   = "Observatory"
	MEASURED_KPI_TYPE      = "Measured"
	QUESTIONNAIRE_KPI_TYPE = "Questionnaire"
)

type KPIService struct {
	Db *gorm.DB
}

func NewKPIService() *KPIService {
	db := database.DB
	err := db.AutoMigrate(&models.Kpis{}, &models.KpiType{}, &models.FeedbackKpi{}, &models.MeasuredKpi{}, &models.ObservatoryKpi{}, &models.QuestionnaireKpi{})
	if err != nil {
		panic(err)
	}
	return &KPIService{Db: db}
}

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
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid feedback KPI data"})
			return
		}

		if feedbackKpi.FeedBack == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid feedback data"})
			return
		}

		kpi, err := controller.CreateKPI(s.Db, kpi)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		feedbackKpi.KpisID = kpi.ID

		result := s.Db.Create(&feedbackKpi)
		if result.Error != nil {
			_ = s.Db.Table("kpis").Delete(&kpi).Error
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, kpi)
	case OBSERVATORY_KPI_TYPE:
		var observatoryKpi models.ObservatoryKpi
		err := c.ShouldBindBodyWith(&observatoryKpi, binding.JSON)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid observatory KPI data"})
			return
		}

		if observatoryKpi.Observatory == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid observatory data"})
			return
		}

		kpi, err := controller.CreateKPI(s.Db, kpi)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		observatoryKpi.KpisID = kpi.ID

		result := s.Db.Create(&observatoryKpi)
		if result.Error != nil {
			_ = s.Db.Table("kpis").Delete(&kpi).Error
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, kpi)
	case MEASURED_KPI_TYPE:
		var measuredKpi models.MeasuredKpi
		err := c.ShouldBindBodyWith(&measuredKpi, binding.JSON)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid measured KPI data"})
			return
		}

		if measuredKpi.Measured == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid measured data"})
			return
		}

		kpi, err := controller.CreateKPI(s.Db, kpi)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		measuredKpi.KpisID = kpi.ID

		result := s.Db.Create(&measuredKpi)
		if result.Error != nil {
			_ = s.Db.Table("kpis").Delete(&kpi).Error
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, kpi)
	case QUESTIONNAIRE_KPI_TYPE:
		var questionnaireKpi models.ReqQuestionnaireKpi
		err := c.ShouldBindBodyWith(&questionnaireKpi, binding.JSON)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid questionnaire KPI data"})
			return
		}

		if questionnaireKpi.Questionnaire == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid questionnaire data"})
			return
		}

		kpi, err := controller.CreateKPI(s.Db, kpi)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for _, q := range questionnaireKpi.Questionnaire {
			question := models.QuestionnaireKpi{
				KpisID:        kpi.ID,
				Questionnaire: q,
			}
			result := s.Db.Table("questionnaire_kpis").Create(&question)
			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
				return
			}
		}

		c.JSON(http.StatusOK, kpi)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid KPI type"})
		return
	}
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch KPIs"})
		return
	}

	c.JSON(http.StatusOK, kpis)
}

func (s *KPIService) GetKPIByID(c *gin.Context) {
	kpiID := c.Param("id")
	var kpi models.Kpis

	if err := controller.GetKPIByID(s.Db, kpiID, &kpi); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "KPI not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch KPI"})
		}
		return
	}

	switch kpi.KpiType {
	case FEEDBACK_KPI_TYPE:
		kpi_data := models.ReqFeedBack{
			ID:         kpi.ID,
			KpiName:    kpi.KpiName,
			KpiType:    kpi.KpiType,
			AssignType: kpi.AssignType,
		}

		var feedbackKpi models.FeedbackKpi
		err := s.Db.Where("kpis_id = ?", kpi.ID).First(&feedbackKpi).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch KPI"})
			return
		}

		kpi_data.Feedback = feedbackKpi.FeedBack

		c.JSON(http.StatusOK, kpi_data)
	case OBSERVATORY_KPI_TYPE:
		kpi_data := models.ReqObservatory{
			ID:         kpi.ID,
			KpiName:    kpi.KpiName,
			KpiType:    kpi.KpiType,
			AssignType: kpi.AssignType,
		}

		var observatoryKpi models.ObservatoryKpi
		err := s.Db.Where("kpis_id = ?", kpi.ID).First(&observatoryKpi).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch KPI"})
			return
		}

		kpi_data.Observatory = observatoryKpi.Observatory

		c.JSON(http.StatusOK, kpi_data)
	case MEASURED_KPI_TYPE:
		kpi_data := models.ReqMeasured{
			ID:         kpi.ID,
			KpiName:    kpi.KpiName,
			KpiType:    kpi.KpiType,
			AssignType: kpi.AssignType,
		}

		var measuredKpi models.MeasuredKpi
		err := s.Db.Where("kpis_id = ?", kpi.ID).First(&measuredKpi).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch KPI"})
			return
		}

		kpi_data.Measured = measuredKpi.Measured

		c.JSON(http.StatusOK, kpi_data)
	case QUESTIONNAIRE_KPI_TYPE:
		kpi_data := models.ReqQuestionnaire{
			ID:         kpi.ID,
			KpiName:    kpi.KpiName,
			KpiType:    kpi.KpiType,
			AssignType: kpi.AssignType,
		}

		var questionnaireKpi []models.QuestionnaireKpi
		err := s.Db.Find(&questionnaireKpi).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch KPI"})
			return
		}

		for _, q := range questionnaireKpi {
			kpi_data.Questionnaire = append(kpi_data.Questionnaire, q.Questionnaire)
		}

		c.JSON(http.StatusOK, kpi_data)
	}
}

func (s *KPIService) UpdateKPI(c *gin.Context) {
	kpiID := c.Param("id")
	var kpi models.Kpis

	if err := s.Db.Where("id = ?", kpiID).First(&kpi).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "KPI not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch KPI"})
		}
		return
	}

	if err := c.ShouldBindJSON(&kpi); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
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

	kpi, err := controller.DeleteKPIByID(s.Db, kpiID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "KPI not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete KPI"})
		}
		return
	}

	switch kpi.KpiType {
	case FEEDBACK_KPI_TYPE:
		err = s.Db.Where("kpis_id = ?", kpi.ID).Delete(&models.FeedbackKpi{}).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete KPI"})
			return
		}
	case OBSERVATORY_KPI_TYPE:
		err = s.Db.Where("kpis_id = ?", kpi.ID).Delete(&models.ObservatoryKpi{}).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete KPI"})
			return
		}
	case MEASURED_KPI_TYPE:
		err = s.Db.Where("kpis_id = ?", kpi.ID).Delete(&models.MeasuredKpi{}).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete KPI"})
			return
		}
	case QUESTIONNAIRE_KPI_TYPE:
		err = s.Db.Where("kpis_id = ?", kpi.ID).Delete(&models.ReqQuestionnaireKpi{}).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete KPI"})
			return
		}
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid KPI type"})
		return
	}

	c.Status(http.StatusNoContent)
}
