package controller

import (
	"github.com/mrehanabbasi/appraisal-system-backend/models"
	"gorm.io/gorm"
)

func CreateAppraisalFlow(db *gorm.DB, appraisalflow models.ApraisalFlow) (models.ApraisalFlow, error) {

	if err := db.Table("apraisal_flows").Create(&appraisalflow).Error; err != nil {
		return appraisalflow, err
	}
	return appraisalflow, nil
}

func GetAppraisalFlowByID(db *gorm.DB, appraisalflow *models.ApraisalFlow, id int) (err error) {
	err = db.Model(&appraisalflow).Preload("FlowSteps").Where("id = ?", id).First(&appraisalflow).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateAppraisalFlow(db *gorm.DB, appraisalflow *models.ApraisalFlow) error {
	if err := db.Table("apraisal_flows").Updates(appraisalflow).Error; err != nil {
		return err
	}
	return nil
}