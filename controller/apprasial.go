package controller

import (
	"fmt"

	log "github.com/mrehanabbasi/appraisal-system-backend/logger"
	"github.com/mrehanabbasi/appraisal-system-backend/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateAppraisal(db *gorm.DB, appraisal models.Apprasial) (models.Apprasial, error) {
	log.Info("Creating appraisal")

	// Check if KPI IDs exist in KPIs table
	for _, kpi := range appraisal.AppraisalKpis {
		var k models.Kpi
		err := db.First(&k, kpi.KpiID).Error
		if err != nil {
			log.Error(err.Error())
			return appraisal, fmt.Errorf("KPI Id not exist")
		}
	}

	if err := db.Create(&appraisal).Error; err != nil {
		log.Error(err.Error())
		return appraisal, err
	}
	return appraisal, nil
}
func GetAppraisalByID(db *gorm.DB, appraisal *models.Apprasial, id int) (err error) {
	log.Info("Getting appraisal by ID")

	err = db.Model(&appraisal).Preload("AppraisalKpis").Where("id = ?", id).First(&appraisal).Error
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func GetAllApprasial(db *gorm.DB, appraisal *[]models.Apprasial) (err error) {
	log.Info("Getting all appraisal")

	err = db.Model(models.Apprasial{}).Preload("AppraisalKpis").Find(&appraisal).Error
	if err != nil {
		return err
	}
	return nil

}

func UpdateAppraisal(db *gorm.DB, appraisal *models.Apprasial, id int) (models.Apprasial, error) {
	log.Info("Updating appraisal")
	// Check if KPI IDs exist in KPIs table
	for _, kpi := range appraisal.AppraisalKpis {
		var k models.Kpi
		err := db.First(&k, kpi.KpiID).Error
		if err != nil {
			log.Error(err.Error())
			return *appraisal, fmt.Errorf("KPI Id not found")
		}
	}

	err := db.Session(&gorm.Session{FullSaveAssociations: true}).Where("id = ?", id).Save(&appraisal).Error
	if err != nil {
		log.Error(err.Error())
		return *appraisal, err
	}
	return *appraisal, nil
}

func DeleteApprasial(db *gorm.DB, appraisal *models.Apprasial, id int) error {
	log.Info("Deleting appraisal")

	err := db.Select(clause.Associations).Delete(&appraisal).Error
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}
