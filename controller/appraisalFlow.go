package controller

import (
	"github.com/mrehanabbasi/appraisal-system-backend/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func GetAllApprasialFlow(flowName, isActive, teamId string, db *gorm.DB, appraisalflow *[]models.ApraisalFlow) (err error) {
	if flowName != "" && isActive != "" && teamId != "" {
		err = db.Model(&appraisalflow).Preload("FlowSteps").Where("flow_name = ? AND is_active = ? AND team_id = ?", flowName, isActive, teamId).Find(&appraisalflow).Error
		if err != nil {
			return err
		}
	} else if flowName != "" {
		err = db.Model(&appraisalflow).Preload("FlowSteps").Where("flow_name LIKE ?", "%"+flowName+"%").Find(&appraisalflow).Error
		if err != nil {
			return err
		}
	} else if flowName != "" && isActive != "" {
		err = db.Model(&appraisalflow).Preload("FlowSteps").Where("flow_name = ? AND is_active = ?", flowName, isActive).Find(&appraisalflow).Error
		if err != nil {
			return err
		}
	} else if flowName != "" && teamId != "" {
		err = db.Model(&appraisalflow).Preload("FlowSteps").Where("flow_name = ? AND team_id = ?", flowName, teamId).Find(&appraisalflow).Error
		if err != nil {
			return err
		}
	} else if isActive != "" && teamId != "" {
		err = db.Model(&appraisalflow).Preload("FlowSteps").Where("is_active = ? AND team_id = ?", isActive, teamId).Find(&appraisalflow).Error
		if err != nil {
			return err
		}
	} else if isActive != "" {
		err = db.Model(&appraisalflow).Preload("FlowSteps").Where("is_active = ?", isActive).Find(&appraisalflow).Error
		if err != nil {
			return err
		}
	} else if teamId != "" {
		err = db.Model(&appraisalflow).Preload("FlowSteps").Where("team_id = ?", teamId).Find(&appraisalflow).Error
		if err != nil {
			return err
		}
	} else {

		err = db.Model(&appraisalflow).Preload("FlowSteps").Find(&appraisalflow).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateAppraisalFlow(db *gorm.DB, appraisalflow *models.ApraisalFlow, id int) error {
	db.Session(&gorm.Session{FullSaveAssociations: true}).Where("id = ?", id).Save(&appraisalflow)
	return nil
}

func DeleteApprasialFlow(db *gorm.DB, appraisalflow *models.ApraisalFlow, id int) error {
	err := db.Select(clause.Associations).Delete(&appraisalflow).Error
	if err != nil {
		return err
	}
	return nil
}
