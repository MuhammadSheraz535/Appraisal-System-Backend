package controller

import (
	"gorm.io/gorm"

	"github.com/mrehanabbasi/appraisal-system-backend/models"
)

func CreateKPI(db *gorm.DB, kpi models.KPI) (models.KPI, error) {
	if err := db.Create(&kpi).Error; err != nil {
		return kpi, err
	}
	return kpi, nil
}

func GetAllKPI(db *gorm.DB, queryParams map[string]string) ([]models.KPI, error) {
	var kpis []models.KPI

	query := db
	if kpiName, ok := queryParams["kpi_name"]; ok {
		query = query.Where("kpi_name = ?", kpiName)
	}

	if assignType, ok := queryParams["assign_type"]; ok {
		query = query.Where("assign_type = ?", assignType)
	}

	if rolesApplicable, ok := queryParams["roles_applicable"]; ok {
		query = query.Where("roles_applicable LIKE ?", "%"+rolesApplicable+"%")
	}

	if err := query.Find(&kpis).Error; err != nil {
		return kpis, err
	}

	return kpis, nil
}

// GetKPI retrieves a single KPI by ID
func GetKPIByID(db *gorm.DB, id uint) (models.KPI, error) {
	var kpi models.KPI

	if err := db.Preload("Measured").Preload("Questionaire").First(&kpi, id).Error; err != nil {
		return kpi, err
	}

	return kpi, nil
}

func UpdateKPI(db *gorm.DB, id uint, updatedKPI models.KPI) (models.KPI, error) {
	var kpi models.KPI

	if err := db.Model(&kpi).Where("id = ?", id).Updates(updatedKPI).Error; err != nil {
		return kpi, err
	}

	if err := db.Model(&models.MeasuredData{}).Where("kpi_id = ?", id).Updates(models.MeasuredData{
		Key:   updatedKPI.Measured.Key,
		Value: updatedKPI.Measured.Value,
	}).Error; err != nil {
		return kpi, err
	}

	if err := db.Model(&models.QuestionaireData{}).Where("kpi_id = ?", id).Updates(models.QuestionaireData{
		Key:   updatedKPI.Questionaire.Key,
		Value: updatedKPI.Questionaire.Value,
	}).Error; err != nil {
		return kpi, err
	}

	if err := db.Preload("Measured").Preload("Questionaire").First(&kpi, id).Error; err != nil {
		return kpi, err
	}

	return kpi, nil
}

func DeleteKPI(db *gorm.DB, id uint) error {
	var kpi models.KPI
	if err := db.First(&kpi, id).Error; err != nil {
		return err
	}

	if err := db.Where("kpi_id = ?", id).Delete(&models.MeasuredData{}, "kpi_id").Error; err != nil {
		return err
	}
	if err := db.Where("kpi_id = ?", id).Delete(&models.QuestionaireData{}, "kpi_id").Error; err != nil {
		return err
	}

	if err := db.Delete(&kpi).Error; err != nil {
		return err
	}

	return nil
}
