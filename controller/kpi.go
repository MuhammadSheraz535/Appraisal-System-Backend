package controller

import (
	"errors"

	log "github.com/mrehanabbasi/appraisal-system-backend/logger"
	"github.com/mrehanabbasi/appraisal-system-backend/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateKPI(db *gorm.DB, kpi *models.Kpi) (*models.Kpi, error) {
	log.Info("Creating new KPI")

	// Check if KPI name already exists
	var count int64
	if err := db.Model(&models.Kpi{}).Where("kpi_name = ?", kpi.KpiName).Count(&count).Error; err != nil {
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

	// Check if KPI exists in the database
	var existingKpi models.Kpi
	if err := db.Model(&models.Kpi{}).First(&existingKpi, kpi.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("kpi with the given id not found")
			return nil, errors.New("kpi not found")
		}
		log.Error(err.Error())
		return nil, err
	}

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

	// Retrieve statements for the existing KPI
	var existingStatements []models.MultiStatementKpiData
	if err := db.Model(&models.MultiStatementKpiData{}).Find(&existingStatements, "kpi_id = ?", kpi.ID).Error; err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// Delete remaining statements if the number of statements is reduced
	if len(existingStatements) > len(kpi.Statements) {
		deletedStatements := existingStatements[len(kpi.Statements):]
		for _, statement := range deletedStatements {
			if err := db.Delete(&statement).Error; err != nil {
				log.Error(err.Error())
				return nil, err
			}
		}
	}

	// Assign statements' IDs to the request statements
	for k := range kpi.Statements {
		if k < len(existingStatements) {
			kpi.Statements[k].ID = existingStatements[k].ID
		}
	}

	// Update KPI record
	if err := db.Session(&gorm.Session{FullSaveAssociations: true}).Where("id = ?", kpi.ID).Save(&kpi).Error; err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return kpi, nil
}

func GetKPIByID(db *gorm.DB, id uint64) (models.Kpi, error) {
	log.Info("Getting KPI by ID")

	var kpi models.Kpi

	if err := db.Model(&models.Kpi{}).Preload("Statements").Where("id = ?", id).First(&kpi).Error; err != nil {
		log.Error(err.Error())
		return kpi, err
	}

	return kpi, nil
}

func GetAllKPI(db *gorm.DB, kpi *[]models.Kpi) error {
	log.Info("Getting all KPIs")

	err := db.Model(&models.Kpi{}).Preload("Statements").Order("id ASC").Find(&kpi).Error
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

	if err := db.Model(&models.Kpi{}).First(&kpi, id).Error; err != nil {
		log.Error("kpi with the given id not found")
		return errors.New("kpi not found")
	}

	// Delete the KPI
	if err := db.Select(clause.Associations).Delete(&kpi).Error; err != nil {
		log.Error("failed to delete kpi")
		return errors.New("failed to delete KPI")
	}

	return nil
}
