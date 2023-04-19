package controller

import (
	"errors"

	log "github.com/mrehanabbasi/appraisal-system-backend/logger"
	"github.com/mrehanabbasi/appraisal-system-backend/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateAppraisalFlow(db *gorm.DB, appraisalFlow models.AppraisalFlow) (models.AppraisalFlow, error) {
	log.Info("Creating appraisal flow")

	var count int64
	if err := db.Table("appraisal_flows").Where("flow_name = ?", appraisalFlow.FlowName).Count(&count).Error; err != nil {
		log.Error(err.Error())
		return appraisalFlow, err
	}
	if count > 0 {
		log.Error("appraisal flow name already exists")
		return appraisalFlow, errors.New("flow name already exists")
	}

	// check uniqueness of step names for all flows
	var stepNames []string
	for _, flow := range appraisalFlow.FlowSteps {
		stepNames = append(stepNames, flow.StepName)
	}
	var stepCount int64
	if err := db.Table("flow_steps").Where("step_name IN (?)", stepNames).Count(&stepCount).Error; err != nil {
		log.Error(err.Error())
		return appraisalFlow, err
	}
	if stepCount > 0 {
		log.Error("flow step name already exists in the same request")
		return appraisalFlow, errors.New("step name already exists")
	}

	if err := db.Create(&appraisalFlow).Error; err != nil {
		log.Error(err.Error())
		return appraisalFlow, err
	}
	return appraisalFlow, nil
}

func GetAppraisalFlowByID(db *gorm.DB, appraisalFlow *models.AppraisalFlow, id int) (err error) {
	log.Info("Getting appraisal flow by ID")

	err = db.Model(&appraisalFlow).Preload("FlowSteps").Where("id = ?", id).First(&appraisalFlow).Error
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func GetAllApprasialFlow(flowName, isActive, teamId string, db *gorm.DB, appraisalFlow *[]models.AppraisalFlow) (err error) {
	log.Info("Getting all appraisal flows")

	if flowName != "" && isActive != "" && teamId != "" {
		err = db.Model(&appraisalFlow).Preload("FlowSteps").Where("flow_name = ? AND is_active = ? AND team_id = ?", flowName, isActive, teamId).Find(&appraisalFlow).Error
		if err != nil {
			log.Error(err.Error())
			return err
		}
	} else if flowName != "" {
		err = db.Model(&appraisalFlow).Preload("FlowSteps").Where("flow_name LIKE ?", "%"+flowName+"%").Find(&appraisalFlow).Error
		if err != nil {
			log.Error(err.Error())
			return err
		}
	} else if flowName != "" && isActive != "" {
		err = db.Model(&appraisalFlow).Preload("FlowSteps").Where("flow_name = ? AND is_active = ?", flowName, isActive).Find(&appraisalFlow).Error
		if err != nil {
			log.Error(err.Error())
			return err
		}
	} else if flowName != "" && teamId != "" {
		err = db.Model(&appraisalFlow).Preload("FlowSteps").Where("flow_name = ? AND team_id = ?", flowName, teamId).Find(&appraisalFlow).Error
		if err != nil {
			log.Error(err.Error())
			return err
		}
	} else if isActive != "" && teamId != "" {
		err = db.Model(&appraisalFlow).Preload("FlowSteps").Where("is_active = ? AND team_id = ?", isActive, teamId).Find(&appraisalFlow).Error
		if err != nil {
			log.Error(err.Error())
			return err
		}
	} else if isActive != "" {
		err = db.Model(&appraisalFlow).Preload("FlowSteps").Where("is_active = ?", isActive).Find(&appraisalFlow).Error
		if err != nil {
			log.Error(err.Error())
			return err
		}
	} else if teamId != "" {
		err = db.Model(&appraisalFlow).Preload("FlowSteps").Where("team_id = ?", teamId).Find(&appraisalFlow).Error
		if err != nil {
			log.Error(err.Error())
			return err
		}
	} else {

		err = db.Model(&appraisalFlow).Preload("FlowSteps").Find(&appraisalFlow).Error
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}

	return nil
}

func UpdateAppraisalFlow(db *gorm.DB, appraisalFlow *models.AppraisalFlow, id int) error {
	log.Info("Updating appraisal flow")

	err := db.Session(&gorm.Session{FullSaveAssociations: true}).Where("id = ?", id).Save(&appraisalFlow).Error
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func DeleteApprasialFlow(db *gorm.DB, appraisalFlow *models.AppraisalFlow, id int) error {
	log.Info("Deleting appraisal flow")

	err := db.Select(clause.Associations).Delete(&appraisalFlow).Error
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}
