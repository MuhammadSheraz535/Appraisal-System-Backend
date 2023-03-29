package controller

import (
	"errors"

	"github.com/mrehanabbasi/appraisal-system-backend/models"
	"gorm.io/gorm"
)

func CreateKPI(db *gorm.DB, kpi models.Kpis) (models.Kpis, error) {
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

func GetAllKPI(db *gorm.DB, kpis *[]models.Kpis) (err error) {
	err = db.Table("kpis").Find(&kpis).Error
	if err != nil {
		return err
	}
	return nil
}

func GetKPIByID(db *gorm.DB, kpiID string, kpi *models.Kpis) error {
	err := db.Table("kpis").Where("id = ?", kpiID).First(kpi).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateKPI(db *gorm.DB, kpi models.Kpis) (models.Kpis, error) {
	// Check if a KPI with the same name already exists in the database
	var existingKpi models.Kpis
	if err := db.Where("kpi_name = ?", kpi.KpiName).First(&existingKpi).Error; err == nil && existingKpi.ID != kpi.ID {
		return models.Kpis{}, errors.New("a KPI with the same name already exists in the database")
	}

	// Update the KPI
	if err := db.Save(&kpi).Error; err != nil {
		return models.Kpis{}, err
	}

	return kpi, nil
}


func DeleteKPIByID(db *gorm.DB, id string) error {
	var kpi models.Kpis
	err := db.Table("kpis").Where("id = ?", id).First(&kpi).Error
	if err != nil {
		return err
	}
	err = db.Delete(&kpi).Error
	if err != nil {
		return err
	}
	return nil
}
