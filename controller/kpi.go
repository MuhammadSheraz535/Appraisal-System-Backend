package controller

import (
	"errors"

	log "github.com/mrehanabbasi/appraisal-system-backend/logger"
	"github.com/mrehanabbasi/appraisal-system-backend/models"
	"gorm.io/gorm"
)

func CreateKPI(db *gorm.DB, kpi *models.Kpi) (*models.Kpi, error) {
	log.Info("Creating new KPI")

	// Check if KPI name already exists
	var count int64
	if err := db.Model(&models.Kpi{}).Where("kpi_name = ? AND id != ?", kpi.KpiName, kpi.ID).Count(&count).Error; err != nil {
		log.Error(err.Error())
		return nil, err
	}
	if count > 0 {
		log.Error("kpi name already exists")
		return nil, errors.New("kpi name already exists")
	}

	// Create new KPI record
	if err := db.Create(kpi).Error; err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return kpi, nil
}

func UpdateKPI(db *gorm.DB, kpi *models.Kpi) (*models.Kpi, error) {
	log.Info("Updating KPI")

	// Check if KPI name already exists
	var count int64
	if err := db.Model(&models.Kpi{}).Where("kpi_name = ? AND id != ?", kpi.KpiName, kpi.ID).Count(&count).Error; err != nil {
		log.Error(err.Error())
		return nil, err
	}
	if count > 0 {
		log.Error("invalid kpi id or kpi name already exists")
		return nil, errors.New("invalid kpi id or kpi name already exists")
	}

	if err := db.First(&kpi, kpi.ID).Error; err != nil {
		log.Error("kpi with the given id not found")
		return nil, errors.New("kpi not found")
	}

	// TODO: Ensure that data from multi_statement_kpi_data table is only updated instead of being deleted

	tx := db.Begin()

	if err := tx.Delete(&models.MultiStatementKpiData{}, "kpi_id = ?", kpi.ID).Error; err != nil {
		tx.Rollback()
		log.Error(err.Error())
		return nil, err
	}

	// Update KPI record
	if err := tx.Save(kpi).Error; err != nil {
		tx.Rollback()
		log.Error(err.Error())
		return nil, err
	}

	tx.Commit()

	return kpi, nil
}

func GetKPIByID(db *gorm.DB, id uint64) (models.Kpi, error) {
	log.Info("Getting KPI by ID")

	var kpi models.Kpi

	if err := db.Where("id = ?", id).First(&kpi).Error; err != nil {
		log.Error(err.Error())
		return kpi, err
	}

	return kpi, nil
}

func GetAllKPI(db *gorm.DB, kpi *[]models.Kpi) (err error) {
	log.Info("Getting all KPIs")

	err = db.Table("kpis").Find(&kpi).Error
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

// DeleteKPI deletes a KPI with the given ID
func DeleteKPI(db *gorm.DB, id uint64) error {
	log.Info("Deleting KPI")

	var kpi models.Kpi

	if err := db.First(&kpi, id).Error; err != nil {
		log.Error("kpi with the given id not found")
		return errors.New("kpi not found")
	}

	// Delete the KPI
	if err := db.Select("Statements").Delete(&kpi).Error; err != nil {
		log.Error("failed to delete kpi")
		return errors.New("failed to delete KPI")
	}

	return nil
}
