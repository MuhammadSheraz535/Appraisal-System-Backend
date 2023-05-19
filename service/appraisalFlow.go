package service

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mrehanabbasi/appraisal-system-backend/constants"
	"github.com/mrehanabbasi/appraisal-system-backend/controller"
	"github.com/mrehanabbasi/appraisal-system-backend/database"
	log "github.com/mrehanabbasi/appraisal-system-backend/logger"
	"github.com/mrehanabbasi/appraisal-system-backend/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AppraisalFlowService struct {
	Db *gorm.DB
}

func NewAppraisalFlowService() *AppraisalFlowService {
	db := database.DB
	err := db.AutoMigrate(models.AppraisalType{}, &models.AppraisalFlow{}, &models.FlowStep{})
	if err != nil {
		log.Panic(err.Error())
		panic(err)
	}

	// Populate appraisal_types table
	err = populateAppraisalTypeTable(db)
	if err != nil {
		panic(err)
	}

	return &AppraisalFlowService{Db: db}
}

func populateAppraisalTypeTable(db *gorm.DB) error {
	appraisalTypes := []string{
		constants.MID_YEAR_APPRAISAL,
		constants.ANNUAL_APPRAISAL,
	}

	appraisalTypesSlice := make([]models.AppraisalType, len(appraisalTypes))

	for k, v := range appraisalTypes {
		newAppraisalType := models.AppraisalType{
			AppraisalType: v,
		}
		if k == 0 {
			newAppraisalType.AppraisalType = constants.MID_YEAR_APPRAISAL
		} else if k == 2 {
			newAppraisalType.AppraisalType = constants.ANNUAL_APPRAISAL
		}

		appraisalTypesSlice[k] = newAppraisalType
	}

	err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&appraisalTypesSlice).Error
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

func (r *AppraisalFlowService) CreateAppraisalFlow(c *gin.Context) {
	log.Info("Initializing CreateAppraisalFlow handler function...")

	var appraisalFlow models.AppraisalFlow
	err := c.ShouldBindJSON(&appraisalFlow)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate AppraisalFlow struct
	err = appraisalFlow.Validate()
	if err != nil {
		errs, ok := controller.ErrValidationSlice(err)
		if !ok {
			log.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Error(err.Error())
		if len(errs) > 1 {
			c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": errs[0]})
		}
		return
	}

	appraisalFlow.ID = 0

	// Validate each FlowStep struct
	for _, flowStep := range appraisalFlow.FlowSteps {
		err = flowStep.Validate()
		if err != nil {
			errs, ok := controller.ErrValidationSlice(err)
			if !ok {
				log.Error(err.Error())
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			log.Error(err.Error())
			if len(errs) > 1 {
				c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": errs[0]})
			}
			return
		}
	}
	// check appraisal type exists
	err = checkAppraisalType(r.Db, appraisalFlow.AppraisalTypeStr)
	if err != nil {
		log.Error("invalid appraisal type")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid appraisal type"})
		return
	}
	dbAppraisalFlow, err := controller.CreateAppraisalFlow(r.Db, &appraisalFlow)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dbAppraisalFlow)
}

func (r *AppraisalFlowService) GetAppraisalFlowByID(c *gin.Context) {
	log.Info("Initializing GetAppraisalFlowByID handler function...")

	id, _ := strconv.ParseUint(c.Param("id"), 0, 64)

	var appraisalFlow models.AppraisalFlow
	err := controller.GetAppraisalFlowByID(r.Db, &appraisalFlow, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("appraisal flow record not found against the given id")
			c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
			return
		}

		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appraisalFlow)
}

func (r *AppraisalFlowService) GetAllAppraisalFlows(c *gin.Context) {
	log.Info("Initializing GetAllAppraisalFlow handler function...")

	var appraisalFlows []models.AppraisalFlow

	flowName := c.Query("flow_name")
	isActive := c.Query("is_active")
	teamId := c.Query("team_id")

	err := controller.GetAllAppraisalFlow(flowName, isActive, teamId, r.Db, &appraisalFlows)

	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appraisalFlows)
}

func (r *AppraisalFlowService) UpdateAppraisalFlow(c *gin.Context) {
	log.Info("Initializing UpdateAppraisalFlow handler function...")

	var appraisalFlow models.AppraisalFlow
	id, _ := strconv.ParseUint(c.Param("id"), 0, 64)

	err := c.ShouldBindJSON(&appraisalFlow)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate AppraisalFlow struct
	err = appraisalFlow.Validate()
	if err != nil {
		errs, ok := controller.ErrValidationSlice(err)
		if !ok {
			log.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Error(err.Error())
		if len(errs) > 1 {
			c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": errs[0]})
		}
		return
	}

	appraisalFlow.ID = id

	// Validate each FlowStep struct
	for _, flowStep := range appraisalFlow.FlowSteps {
		err = flowStep.Validate()
		if err != nil {
			errs, ok := controller.ErrValidationSlice(err)
			if !ok {
				log.Error(err.Error())
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			log.Error(err.Error())
			if len(errs) > 1 {
				c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": errs[0]})
			}
			return
		}
	}

	// check appraisal type exists
	err = checkAppraisalType(r.Db, appraisalFlow.AppraisalTypeStr)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid appraisal type"})
		return
	}
	appraisalFlow.ID = id

	// calling controller update method
	err = controller.UpdateAppraisalFlow(r.Db, &appraisalFlow)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appraisalFlow)
}

func (r *AppraisalFlowService) DeleteAppraisalFlow(c *gin.Context) {
	log.Info("Initializing DeleteAppraisalFlow handler function...")

	var appraisalFlow models.AppraisalFlow
	id, _ := strconv.ParseUint(c.Param("id"), 0, 64)
	appraisalFlow.ID = id

	err := controller.GetAppraisalFlowByID(r.Db, &appraisalFlow, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("appraisal flow record not found against the given id")
			c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
			return
		}

		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = controller.DeleteAppraisalFlow(r.Db, &appraisalFlow, id)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
