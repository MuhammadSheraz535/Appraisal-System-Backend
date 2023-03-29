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
	// Update the KPI
	if err := db.Save(&kpi).Error; err != nil {
		return models.Kpis{}, err
	}

	return kpi, nil
}

func DeleteKPIByID(db *gorm.DB, id string) (*models.Kpis, error) {
	var kpi models.Kpis
	err := GetKPIByID(db, id, &kpi)
	if err != nil {
		return nil, err
	}
	err = db.Delete(&kpi).Error
	if err != nil {
		return nil, err
	}

	return &kpi, nil
}
