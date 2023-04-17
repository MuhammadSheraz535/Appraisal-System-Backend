package controller

import (
	"errors"

	"github.com/mrehanabbasi/appraisal-system-backend/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateAppraisalFlow(db *gorm.DB, appraisalFlow models.AppraisalFlow) (models.AppraisalFlow, error) {
	var count int64
	if err := db.Table("appraisal_flows").Where("flow_name = ?", appraisalFlow.FlowName).Count(&count).Error; err != nil {
		return appraisalFlow, err
	}
	if count > 0 {
		return appraisalFlow, errors.New("flow name already exists")
	}

	// check uniqueness of step names for all flows
	var stepNames []string
	for _, flow := range appraisalFlow.FlowSteps {
		stepNames = append(stepNames, flow.StepName)
	}
	var stepCount int64
	if err := db.Table("flow_steps").Where("step_name IN (?)", stepNames).Count(&stepCount).Error; err != nil {
		return appraisalFlow, err
	}
	if stepCount > 0 {
		return appraisalFlow, errors.New("step name already exists")
	}

	if err := db.Create(&appraisalFlow).Error; err != nil {
		return appraisalFlow, err
	}
	return appraisalFlow, nil
}



func GetAppraisalFlowByID(db *gorm.DB, appraisalFlow *models.AppraisalFlow, id int) (err error) {
	err = db.Model(&appraisalFlow).Preload("FlowSteps").Where("id = ?", id).First(&appraisalFlow).Error
	if err != nil {
		return err
	}
	return nil
}

func GetAllApprasialFlow(flowName, isActive, teamId string, db *gorm.DB, appraisalFlow *[]models.AppraisalFlow) (err error) {
	if flowName != "" && isActive != "" && teamId != "" {
		err = db.Model(&appraisalFlow).Preload("FlowSteps").Where("flow_name = ? AND is_active = ? AND team_id = ?", flowName, isActive, teamId).Find(&appraisalFlow).Error
		if err != nil {
			return err
		}
	} else if flowName != "" {
		err = db.Model(&appraisalFlow).Preload("FlowSteps").Where("flow_name LIKE ?", "%"+flowName+"%").Find(&appraisalFlow).Error
		if err != nil {
			return err
		}
	} else if flowName != "" && isActive != "" {
		err = db.Model(&appraisalFlow).Preload("FlowSteps").Where("flow_name = ? AND is_active = ?", flowName, isActive).Find(&appraisalFlow).Error
		if err != nil {
			return err
		}
	} else if flowName != "" && teamId != "" {
		err = db.Model(&appraisalFlow).Preload("FlowSteps").Where("flow_name = ? AND team_id = ?", flowName, teamId).Find(&appraisalFlow).Error
		if err != nil {
			return err
		}
	} else if isActive != "" && teamId != "" {
		err = db.Model(&appraisalFlow).Preload("FlowSteps").Where("is_active = ? AND team_id = ?", isActive, teamId).Find(&appraisalFlow).Error
		if err != nil {
			return err
		}
	} else if isActive != "" {
		err = db.Model(&appraisalFlow).Preload("FlowSteps").Where("is_active = ?", isActive).Find(&appraisalFlow).Error
		if err != nil {
			return err
		}
	} else if teamId != "" {
		err = db.Model(&appraisalFlow).Preload("FlowSteps").Where("team_id = ?", teamId).Find(&appraisalFlow).Error
		if err != nil {
			return err
		}
	} else {

		err = db.Model(&appraisalFlow).Preload("FlowSteps").Find(&appraisalFlow).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateAppraisalFlow(db *gorm.DB, appraisalFlow *models.AppraisalFlow, id int) error {
	db.Session(&gorm.Session{FullSaveAssociations: true}).Where("id = ?", id).Save(&appraisalFlow)
	return nil
}

func DeleteApprasialFlow(db *gorm.DB, appraisalFlow *models.AppraisalFlow, id int) error {
	err := db.Select(clause.Associations).Delete(&appraisalFlow).Error
	if err != nil {
		return err
	}
	return nil
}
