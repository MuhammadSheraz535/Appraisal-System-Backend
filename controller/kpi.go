package controller

import (
	"errors"
	"fmt"

	"github.com/mrehanabbasi/appraisal-system-backend/models"
	"gorm.io/gorm"
)

func CreateKPI(db *gorm.DB, kpiCommon models.KpiCommon, feedbackKpi models.FeedbackKpi, observatoryKpi models.ObservatoryKpi, questionaireKpi models.QuestionaireKpi, measuredKpi models.MeasuredKpi) (models.KpiCommon, error) {
	if err := db.Create(&kpiCommon).Error; err != nil {
		return kpiCommon, err
	}

	feedbackKpi.KpiCommon = kpiCommon
	if err := db.Create(&feedbackKpi).Error; err != nil {
		return kpiCommon, err
	}

	observatoryKpi.KpiCommon = kpiCommon
	if err := db.Create(&observatoryKpi).Error; err != nil {
		return kpiCommon, err
	}

	questionaireKpi.KpiCommon = kpiCommon
	if err := db.Create(&questionaireKpi).Error; err != nil {
		return kpiCommon, err
	}

	measuredKpi.KpiCommon = kpiCommon
	if err := db.Create(&measuredKpi).Error; err != nil {
		return kpiCommon, err
	}

	return kpiCommon, nil
}

func GetAllKPI(db *gorm.DB, kpiType string, kpiName string, assignType string) (interface{}, error) {
	var kpis interface{}

	switch kpiType {
	case "feedback":
		var feedbackKpis []models.FeedbackKpi
		db := db.Where("kpi_type = ?", kpiType)
		if kpiName != "" {
			db = db.Where("kpi_name = ?", kpiName)
		}
		if assignType != "" {
			db = db.Where("assign_type = ?", assignType)
		}
		if err := db.Find(&feedbackKpis).Error; err != nil {
			return nil, err
		}
		kpis = feedbackKpis
	case "observatory":
		var observatoryKpis []models.ObservatoryKpi
		db := db.Where("kpi_type = ?", kpiType)
		if kpiName != "" {
			db = db.Where("kpi_name = ?", kpiName)
		}
		if assignType != "" {
			db = db.Where("assign_type = ?", assignType)
		}
		if err := db.Find(&observatoryKpis).Error; err != nil {
			return nil, err
		}
		kpis = observatoryKpis
	case "questionnaire":
		var questionnaireKpis []models.QuestionaireKpi
		db := db.Where("kpi_type = ?", kpiType)
		if kpiName != "" {
			db = db.Where("kpi_name = ?", kpiName)
		}
		if assignType != "" {
			db = db.Where("assign_type = ?", assignType)
		}
		if err := db.Find(&questionnaireKpis).Error; err != nil {
			return nil, err
		}
		kpis = questionnaireKpis
	case "measured":
		var measuredKpis []models.MeasuredKpi
		db := db.Where("kpi_type = ?", kpiType)
		if kpiName != "" {
			db = db.Where("kpi_name = ?", kpiName)
		}
		if assignType != "" {
			db = db.Where("assign_type = ?", assignType)
		}
		if err := db.Find(&measuredKpis).Error; err != nil {
			return nil, err
		}
		kpis = measuredKpis
	default:
		return nil, errors.New("Invalid KpiType query parameter")
	}

	return kpis, nil
}

func GetKPIByID(db *gorm.DB, kpiType string, kpiID uint) (interface{}, error) {
	var kpi interface{}

	switch kpiType {
	case "feedback":
		var feedbackKpi models.FeedbackKpi
		if err := db.Where("id = ?", kpiID).First(&feedbackKpi).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("Feedback KPI not found")
			}
			return nil, err
		}
		kpi = feedbackKpi
	case "observatory":
		var observatoryKpi models.ObservatoryKpi
		if err := db.Where("id = ?", kpiID).First(&observatoryKpi).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("Observatory KPI not found")
			}
			return nil, err
		}
		kpi = observatoryKpi
	case "questionnaire":
		var questionnaireKpi models.QuestionaireKpi
		if err := db.Where("id = ?", kpiID).First(&questionnaireKpi).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("Questionnaire KPI not found")
			}
			return nil, err
		}
		kpi = questionnaireKpi
	case "measured":
		var measuredKpi models.MeasuredKpi
		if err := db.Where("id = ?", kpiID).First(&measuredKpi).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("Measured KPI not found")
			}
			return nil, err
		}
		kpi = measuredKpi
	default:
		return nil, errors.New("Invalid KpiType query parameter")
	}

	return kpi, nil
}

func UpdateKPI(db *gorm.DB, kpiType string, kpiID uint, kpi interface{}) error {
	switch kpiType {
	case "feedback":
		feedbackKpi, ok := kpi.(*models.FeedbackKpi)
		if !ok {
			return fmt.Errorf("invalid feedback KPI data")
		}

		if err := db.Model(&models.FeedbackKpi{}).Where("id = ?", kpiID).Updates(feedbackKpi).Error; err != nil {
			return fmt.Errorf("failed to update feedback KPI: %v", err)
		}

	case "observatory":
		observatoryKpi, ok := kpi.(*models.ObservatoryKpi)
		if !ok {
			return fmt.Errorf("invalid observatory KPI data")
		}

		if err := db.Model(&models.ObservatoryKpi{}).Where("id = ?", kpiID).Updates(observatoryKpi).Error; err != nil {
			return fmt.Errorf("failed to update observatory KPI: %v", err)
		}

	case "questionnaire":
		questionnaireKpi, ok := kpi.(*models.QuestionaireKpi)
		if !ok {
			return fmt.Errorf("invalid questionnaire KPI data")
		}

		if err := db.Model(&models.QuestionaireKpi{}).Where("id = ?", kpiID).Updates(questionnaireKpi).Error; err != nil {
			return fmt.Errorf("failed to update questionnaire KPI: %v", err)
		}

	case "measured":
		measuredKpi, ok := kpi.(*models.MeasuredKpi)
		if !ok {
			return fmt.Errorf("invalid measured KPI data")
		}

		if err := db.Model(&models.MeasuredKpi{}).Where("id = ?", kpiID).Updates(measuredKpi).Error; err != nil {
			return fmt.Errorf("failed to update measured KPI: %v", err)
		}

	default:
		return fmt.Errorf("invalid KPI type")
	}

	return nil
}

func DeleteKPI(db *gorm.DB, kpiType string, kpiID uint) error {
	var err error
	switch kpiType {
	case "feedback":
		var feedbackKpi models.FeedbackKpi
		if err = db.Where("id = ?", kpiID).First(&feedbackKpi).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("Feedback KPI not found")
			}
			return err
		}

		if err = db.Delete(&feedbackKpi).Error; err != nil {
			return err
		}

	case "observatory":
		var observatoryKpi models.ObservatoryKpi
		if err = db.Where("id = ?", kpiID).First(&observatoryKpi).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("Observatory KPI not found")
			}
			return err
		}

		if err = db.Delete(&observatoryKpi).Error; err != nil {
			return err
		}

	case "questionnaire":
		var questionnaireKpi models.QuestionaireKpi
		if err = db.Where("id = ?", kpiID).First(&questionnaireKpi).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("Questionnaire KPI not found")
			}
			return err
		}

		if err = db.Delete(&questionnaireKpi).Error; err != nil {
			return err
		}

	case "measured":
		var measuredKpi models.MeasuredKpi
		if err = db.Where("id = ?", kpiID).First(&measuredKpi).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("Measured KPI not found")
			}
			return err
		}

		if err = db.Delete(&measuredKpi).Error; err != nil {
			return err
		}

	default:
		return fmt.Errorf("Invalid KPI type")
	}

	return nil
}
