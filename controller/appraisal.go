package controller

import (
	"errors"

	log "github.com/mrehanabbasi/appraisal-system-backend/logger"
	"github.com/mrehanabbasi/appraisal-system-backend/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateAppraisal(db *gorm.DB, appraisal *models.Appraisal) (*models.Appraisal, error) {
	log.Info("Creating appraisal")

	err := checkKpiIdsExist(db, &appraisal.AppraisalKpis)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// Check if appraisal name already exists
	var count int64
	if err := db.Model(&models.Appraisal{}).Where("appraisal_name = ?", appraisal.AppraisalName).Count(&count).Error; err != nil {
		log.Error(err.Error())
		return nil, err
	}
	if count > 0 {
		log.Error("appraisal name already exists")
		return nil, errors.New("appraisal name already exists")
	}

	if err := db.Create(&appraisal).Error; err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return appraisal, nil
}

func GetAppraisalByID(db *gorm.DB, appraisal *models.Appraisal, id uint64) error {
	log.Info("Getting appraisal by ID")

	err := db.Model(&models.Appraisal{}).Preload("AppraisalKpis").Where("id = ?", id).First(&appraisal).Error
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

func GetAllAppraisals(db *gorm.DB, appraisal *[]models.Appraisal) (err error) {
	log.Info("Getting all appraisals")

	err = db.Model(models.Appraisal{}).Preload("AppraisalKpis").Find(&appraisal).Error
	if err != nil {
		return err
	}

	return nil
}

func UpdateAppraisal(db *gorm.DB, appraisal *models.Appraisal, id uint64) (*models.Appraisal, error) {
	log.Info("Updating appraisal")

	err := checkKpiIdsExist(db, &appraisal.AppraisalKpis)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// Check if appraisal name already exists
	var count int64
	if err := db.Model(&models.Appraisal{}).Where("appraisal_name = ? AND id != ?", appraisal.AppraisalName, appraisal.ID).Count(&count).Error; err != nil {
		log.Error(err.Error())
		return nil, err
	}
	if count > 0 {
		log.Error("appraisal name already exists")
		return nil, errors.New("appraisal name already exists")
	}

	err = db.Session(&gorm.Session{FullSaveAssociations: true}).Where("id = ?", id).Save(&appraisal).Error
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return appraisal, nil
}

func DeleteAppraisal(db *gorm.DB, appraisal *models.Appraisal, id uint64) error {
	log.Info("Deleting appraisal")

	err := db.Select(clause.Associations).Delete(&appraisal).Error
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

// Checks if KPI IDs exist in kpis table
func checkKpiIdsExist(db *gorm.DB, appraisalKpis *[]models.AppraisalKpi) error {
	var kpiIds []uint64
	for _, aKpi := range *appraisalKpis {
		kpiIds = append(kpiIds, aKpi.KpiID)
	}

	var kpis []models.Kpi
	if err := db.Model(&models.Kpi{}).Find(&kpis, kpiIds).Error; err != nil {
		log.Error(err.Error())
		return errors.New("kpi id does not exist")
	}

	return nil
}
