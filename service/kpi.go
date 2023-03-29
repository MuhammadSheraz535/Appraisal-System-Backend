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

const (
	ASSIGN_TYPE_ROLE       = "Role"
	ASSIGN_TYPE_TEAM       = "Team"
	ASSIGN_TYPE_INDIVIDUAL = "Individual"
)

type KPIService struct {
	Db *gorm.DB
}

func NewKPIService() *KPIService {
	db := database.DB
	err := db.AutoMigrate(&models.Kpis{}, &models.KpiType{}, &models.AssignType{}, &models.FeedbackKpi{}, &models.MeasuredKpi{}, &models.ObservatoryKpi{}, &models.QuestionnaireKpi{})
	if err != nil {
		panic(err)
	}

	// Populate assign_types table
	err = populateAssignTypeTable(db)
	if err != nil {
		panic(err)
	}

	// Populate kpi_types table
	err = populateKpiTypeTable(db)
	if err != nil {
		panic(err)
	}

	return &KPIService{Db: db}
}

func populateKpiTypeTable(db *gorm.DB) error {
	kpiTypes := []string{
		FEEDBACK_KPI_TYPE,
		OBSERVATORY_KPI_TYPE,
		MEASURED_KPI_TYPE,
		QUESTIONNAIRE_KPI_TYPE,
	}

	for _, k := range kpiTypes {
		newKpiType := models.KpiType{
			KpiType: k,
		}
		err := db.Create(&newKpiType).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func populateAssignTypeTable(db *gorm.DB) error {
	assignTypes := []string{
		ASSIGN_TYPE_ROLE,
		ASSIGN_TYPE_TEAM,
		ASSIGN_TYPE_INDIVIDUAL,
	}

	for i, a := range assignTypes {
		newAssignType := models.AssignType{
			AssignTypeId: uint(i),
			AssignType:   a,
		}
		err := db.Create(&newAssignType).Error
		if err != nil {
			return err
		}
	}

	return nil
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
				_ = s.Db.Table("kpis").Delete(&kpi).Error
				c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
				return
			}
		}

		c.JSON(http.StatusCreated, kpi)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid KPI type"})
		return
	}
}

func (s *KPIService) GetAllKPI(c *gin.Context) {
	var kpis []models.Kpis
	db := s.Db

	kpiName := c.Query("kpi_name")
	assignType := c.Query("assign_type")
	kpiType := c.Query("kpi_type")

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

	var allKpis []interface{}

	for _, kpi := range kpis {
		switch kpi.KpiType {
		case FEEDBACK_KPI_TYPE:
			kpi_data := models.ReqFeedBack{
				ID:            kpi.ID,
				KpiName:       kpi.KpiName,
				KpiType:       kpi.KpiType,
				AssignType:    kpi.AssignType,
				ApplicableFor: kpi.ApplicableFor,
			}

			var feedbackKpi models.FeedbackKpi
			err := db.Where("kpis_id = ?", kpi.ID).First(&feedbackKpi).Error
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch KPI"})
				return
			}

			kpi_data.Feedback = feedbackKpi.FeedBack

			allKpis = append(allKpis, kpi_data)
		case OBSERVATORY_KPI_TYPE:
			kpi_data := models.ReqObservatory{
				ID:            kpi.ID,
				KpiName:       kpi.KpiName,
				KpiType:       kpi.KpiType,
				AssignType:    kpi.AssignType,
				ApplicableFor: kpi.ApplicableFor,
			}

			var observatoryKpi models.ObservatoryKpi
			err := s.Db.Where("kpis_id = ?", kpi.ID).First(&observatoryKpi).Error
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch KPI"})
				return
			}

			kpi_data.Observatory = observatoryKpi.Observatory

			allKpis = append(allKpis, kpi_data)
		case MEASURED_KPI_TYPE:
			kpi_data := models.ReqMeasured{
				ID:            kpi.ID,
				KpiName:       kpi.KpiName,
				KpiType:       kpi.KpiType,
				AssignType:    kpi.AssignType,
				ApplicableFor: kpi.ApplicableFor,
			}

			var measuredKpi models.MeasuredKpi
			err := s.Db.Where("kpis_id = ?", kpi.ID).First(&measuredKpi).Error
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch KPI"})
				return
			}

			kpi_data.Measured = measuredKpi.Measured

			allKpis = append(allKpis, kpi_data)
		case QUESTIONNAIRE_KPI_TYPE:
			kpi_data := models.ReqQuestionnaire{
				ID:            kpi.ID,
				KpiName:       kpi.KpiName,
				KpiType:       kpi.KpiType,
				AssignType:    kpi.AssignType,
				ApplicableFor: kpi.ApplicableFor,
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

			allKpis = append(allKpis, kpi_data)
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid KPI type"})
			return
		}
	}

	c.JSON(http.StatusOK, allKpis)
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
			ID:            kpi.ID,
			KpiName:       kpi.KpiName,
			KpiType:       kpi.KpiType,
			AssignType:    kpi.AssignType,
			ApplicableFor: kpi.ApplicableFor,
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
			ID:            kpi.ID,
			KpiName:       kpi.KpiName,
			KpiType:       kpi.KpiType,
			AssignType:    kpi.AssignType,
			ApplicableFor: kpi.ApplicableFor,
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
			ID:            kpi.ID,
			KpiName:       kpi.KpiName,
			KpiType:       kpi.KpiType,
			AssignType:    kpi.AssignType,
			ApplicableFor: kpi.ApplicableFor,
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
			ID:            kpi.ID,
			KpiName:       kpi.KpiName,
			KpiType:       kpi.KpiType,
			AssignType:    kpi.AssignType,
			ApplicableFor: kpi.ApplicableFor,
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
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid KPI type"})
		return
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

	if err := c.ShouldBindBodyWith(&kpi, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

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

		kpi, err := controller.UpdateKPI(s.Db, kpi)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = s.Db.Where("kpis_id = ?", kpi.ID).Save(&feedbackKpi).Error
		if err != nil {
			// TODO: Get old kpi and replace it with new one in case of error
			// _ = s.Db.Table("kpis").Delete(&kpi).Error
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

		kpi, err := controller.UpdateKPI(s.Db, kpi)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		result := s.Db.Where("kpis_id = ?", kpi.ID).Save(&observatoryKpi)
		if result.Error != nil {
			// TODO: Get old kpi and replace it with new one in case of error
			// _ = s.Db.Table("kpis").Delete(&kpi).Error
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

		kpi, err := controller.UpdateKPI(s.Db, kpi)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		result := s.Db.Where("kpis_id = ?", kpi.ID).Save(&measuredKpi)
		if result.Error != nil {
			// TODO: Get old kpi and replace it with new one in case of error
			// _ = s.Db.Table("kpis").Delete(&kpi).Error
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

		kpi, err := controller.UpdateKPI(s.Db, kpi)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		_ = s.Db.Where("kpis_id = ?", kpi.ID).Delete(&models.QuestionnaireKpi{}).Error
		for _, q := range questionnaireKpi.Questionnaire {
			question := models.QuestionnaireKpi{
				KpisID:        kpi.ID,
				Questionnaire: q,
			}
			result := s.Db.Table("questionnaire_kpis").Where("kpis_id = ?", kpi.ID).Save(&question)
			if result.Error != nil {
				// TODO: Get old kpi and replace it with new one in case of error
				// _ = s.Db.Table("kpis").Delete(&kpi).Error
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
		err = s.Db.Where("kpis_id = ?", kpi.ID).Delete(&models.QuestionnaireKpi{}).Error
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
