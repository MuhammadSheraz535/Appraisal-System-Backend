package controller

import (
	"errors"

	"github.com/mrehanabbasi/appraisal-system-backend/models"
	"gorm.io/gorm"
)

func CreateKPI(db *gorm.DB, kpi models.Kpi) (models.Kpi, error) {
	// Check if KPI name already exists
	var count int64
	if err := db.Table("kpis").Where("kpi_name = ?", kpi.KpiName).Count(&count).Error; err != nil {
		return kpi, err
	}
	if count > 0 {
		return kpi, errors.New("KPI name already exists")
	}
	// Create new KPI record
	if err := db.Create(&kpi).Error; err != nil {
		return kpi, err
	}
	return kpi, nil
}

func UpdateKPI(db *gorm.DB, kpi models.Kpi) (models.Kpi, error) {
	// Check if KPI name already exists
	var count int64
	if err := db.Table("kpis").Where("kpi_name = ? AND id != ?", kpi.KpiName, kpi.ID).Count(&count).Error; err != nil {
		return kpi, err
	}
	if count > 0 {
		return kpi, errors.New("KPI name already exists")
	}
	// Update KPI record
	if err := db.Save(&kpi).Error; err != nil {
		return kpi, err
	}
	return kpi, nil
}

func GetKPIByID(db *gorm.DB, id uint) (models.Kpi, error) {
	var kpi models.Kpi

	if err := db.Where("id = ?", id).First(&kpi).Error; err != nil {
		return kpi, err
	}

	return kpi, nil
}

func GetAllKPI(db *gorm.DB, kpi *[]models.Kpi) (err error) {
	err = db.Table("kpis").Find(&kpi).Error
	if err != nil {
		return err
	}
	return nil
}

// DeleteKPI deletes a KPI with the given ID
func DeleteKPI(db *gorm.DB, id string) error {
	var kpi models.Kpi

	if err := db.First(&kpi, id).Error; err != nil {
		return errors.New("KPI not found")
	}

	tx := db.Begin()

	// Delete associated statements if KPI type is MultiKpi
	if kpi.KpiType == "Measured" || kpi.KpiType == "Questionnaire" {
		if err := tx.Where("kpi_id = ?", kpi.ID).Delete(&models.MultiStatementKpiData{}).Error; err != nil {
			tx.Rollback()
			return errors.New("failed to delete KPI statements")
		}
	}

	// Delete the KPI
	if err := tx.Delete(&kpi).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to delete KPI")
	}

	if err := tx.Commit().Error; err != nil {
		return errors.New("failed to commit transaction")
	}

	return nil
}
