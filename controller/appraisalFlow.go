package controller

import (
	"errors"

	log "github.com/mrehanabbasi/appraisal-system-backend/logger"
	"github.com/mrehanabbasi/appraisal-system-backend/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateAppraisalFlow(db *gorm.DB, appraisalFlow *models.AppraisalFlow) (*models.AppraisalFlow, error) {
	log.Info("Creating appraisal flow")
	//check flow name exists in database
	var count int64
	if err := db.Model(&models.AppraisalFlow{}).Where("flow_name = ?", appraisalFlow.FlowName).Count(&count).Error; err != nil {
		log.Error(err.Error())
		return nil, err
	}
	if count > 0 {
		log.Error("appraisal flow name already exists")
		return nil, errors.New("flow name already exists")
	}

	// check uniqueness of step names for current request
	var stepNames []string
	for _, flow := range appraisalFlow.FlowSteps {
		if contains(stepNames, flow.StepName) {
			log.Error("flow step name already exists in the same request")
			return nil, errors.New("step name already exists")
		}
		stepNames = append(stepNames, flow.StepName)
	}

	if err := db.Create(&appraisalFlow).Error; err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return appraisalFlow, nil
}

func GetAppraisalFlowByID(db *gorm.DB, appraisalFlow *models.AppraisalFlow, id uint64) error {
	log.Info("Getting appraisal flow by ID")

	err := db.Model(&models.AppraisalFlow{}).Preload("FlowSteps").Where("id = ?", id).First(&appraisalFlow).Error
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func GetAllAppraisalFlow(flowName, isActive, teamId string, db *gorm.DB, appraisalFlows *[]models.AppraisalFlow) error {
	log.Info("Getting all appraisal flows")

	db = db.Model(&models.AppraisalFlow{}).Preload("FlowSteps")

	if flowName != "" {
		db = db.Where("flow_name LIKE ?", "%"+flowName+"%")
	}

	if isActive != "" {
		db = db.Where("is_active = ?", isActive)
	}

	if teamId != "" {
		db = db.Where("team_id = ?", teamId)
	}

	err := db.Order("id ASC").Find(&appraisalFlows).Error
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

func UpdateAppraisalFlow(db *gorm.DB, appraisalFlow *models.AppraisalFlow) error {
	log.Info("Updating appraisal flow")

	// Check if Appraisal exists in the database
	var existingAppraisalFlow models.AppraisalFlow
	if err := db.Model(&models.AppraisalFlow{}).First(&existingAppraisalFlow, appraisalFlow.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("Appraisal flow with the given id not found")
			return errors.New("appraisal flow not found")
		}
		log.Error(err.Error())
		return err
	}

	//check flow name exists in database
	var count int64
	if err := db.Model(&models.AppraisalFlow{}).Where("flow_name = ? AND id != ?", appraisalFlow.FlowName, appraisalFlow.ID).Count(&count).Error; err != nil {
		log.Error(err.Error())
		return err
	}
	if count > 0 {
		log.Error("appraisal flow name already exists")
		return errors.New("flow name already exists")
	}

	// check uniqueness of step names for current request
	var stepNames []string
	for _, flow := range appraisalFlow.FlowSteps {
		if contains(stepNames, flow.StepName) {
			log.Error("flow step name already exists in the same request")
			return errors.New("step name already exists")
		}
		stepNames = append(stepNames, flow.StepName)
	}

	// Retrieve flow steps for the existing AppraisalFlow
	var existingFlowSteps []models.FlowStep
	if err := db.Model(&models.FlowStep{}).Find(&existingFlowSteps, "flow_id = ?", appraisalFlow.ID).Error; err != nil {
		log.Error(err.Error())
		return err
	}

	// Delete remaining flow steps if the number of steps is reduced
	if len(existingFlowSteps) > len(appraisalFlow.FlowSteps) {
		deletedFlowSteps := existingFlowSteps[len(appraisalFlow.FlowSteps):]
		for _, flowStep := range deletedFlowSteps {
			if err := db.Delete(&flowStep).Error; err != nil {
				log.Error(err.Error())
				return err
			}
		}
	}

	// Assign flow steps' IDs to the request flow steps
	for k := range appraisalFlow.FlowSteps {
		if k < len(existingFlowSteps) {
			appraisalFlow.FlowSteps[k].ID = existingFlowSteps[k].ID
		}
	}

	err := db.Session(&gorm.Session{FullSaveAssociations: true}).Where("id = ?", appraisalFlow.ID).Save(&appraisalFlow).Error
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

func DeleteAppraisalFlow(db *gorm.DB, appraisalFlow *models.AppraisalFlow, id uint64) error {
	log.Info("Deleting appraisal flow")

	err := db.Select(clause.Associations).Delete(&appraisalFlow).Error
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

// Contains tells whether a contains x.
func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
