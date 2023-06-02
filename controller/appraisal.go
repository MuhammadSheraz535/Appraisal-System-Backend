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

	// Check if KPI IDs exist in KPIs table
	for _, kpi := range appraisal.AppraisalKpis {
		var k models.Kpi
		err := db.Model(&models.Kpi{}).First(&k, kpi.KpiID).Error
		if err != nil {
			log.Error(err.Error())
			return appraisal, errors.New("kpi id does not exist")
		}
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

	err := db.Model(&models.Appraisal{}).Preload("AppraisalKpis").Preload("EmployeesList").Where("id = ?", id).First(&appraisal).Error
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

func GetAllAppraisals(db *gorm.DB, appraisal *[]models.Appraisal) (err error) {
	log.Info("Getting all appraisals")

	err = db.Model(models.Appraisal{}).Preload("AppraisalKpis").Preload("EmployeesList").Order("id ASC").Find(&appraisal).Error
	if err != nil {
		return err
	}

	return nil
}

func UpdateAppraisal(db *gorm.DB, appraisal *models.Appraisal) (*models.Appraisal, error) {
	log.Info("Updating appraisal")

	// Check if appraisal exists in the database
	var existingAppraisal models.Appraisal
	if err := db.Model(&models.Appraisal{}).First(&existingAppraisal, appraisal.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("appraisal with the given id not found")
			return nil, errors.New("appraisal not found")
		}
		log.Error(err.Error())
		return nil, err
	}
	// Check if KPI IDs exist in KPIs table
	for _, kpi := range appraisal.AppraisalKpis {
		var k models.Kpi
		err := db.Model(&models.Kpi{}).First(&k, kpi.KpiID).Error
		if err != nil {
			log.Error(err.Error())
			return appraisal, errors.New("kpi id does not exist")
		}
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
	// Retrieve AppraisalKpis for the existing Appraisal
	var existingAppraisalKpis []models.AppraisalKpi
	if err := db.Model(&models.AppraisalKpi{}).Find(&existingAppraisalKpis, "appraisal_id = ?", appraisal.ID).Error; err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// Delete remaining AppraisalKpis if the number of KPIs is reduced
	if len(existingAppraisalKpis) > len(appraisal.AppraisalKpis) {
		deletedAppraisalKpis := existingAppraisalKpis[len(appraisal.AppraisalKpis):]
		for _, appraisalKpi := range deletedAppraisalKpis {
			if err := db.Delete(&appraisalKpi).Error; err != nil {
				log.Error(err.Error())
				return nil, err
			}
		}
	}

	// Assign AppraisalKpis' IDs to the request AppraisalKpis
	for k := range appraisal.AppraisalKpis {
		if k < len(existingAppraisalKpis) {
			appraisal.AppraisalKpis[k].ID = existingAppraisalKpis[k].ID
		}
	}
	err := db.Session(&gorm.Session{FullSaveAssociations: true}).Where("id = ?", appraisal.ID).Save(&appraisal).Error
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
