package service

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mrehanabbasi/appraisal-system-backend/controller"
	"github.com/mrehanabbasi/appraisal-system-backend/database"
	log "github.com/mrehanabbasi/appraisal-system-backend/logger"
	"github.com/mrehanabbasi/appraisal-system-backend/models"
	"gorm.io/gorm"
)

type AppraisalService struct {
	Db *gorm.DB
}

func NewAppraisalService() *AppraisalService {
	db := database.DB
	// TODO: Remove migration of Score to it's own service
	err := db.AutoMigrate(&models.Appraisal{}, models.AppraisalKpi{}, models.Score{})
	if err != nil {
		panic(err)
	}

	return &AppraisalService{Db: db}
}

func (r *AppraisalService) CreateAppraisal(c *gin.Context) {
	log.Info("Initializing CreateAppraisal handler function...")

	var appraisal models.Appraisal
	err := c.ShouldBindJSON(&appraisal)
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

	// Validate each FlowStep struct
	for _, ak := range appraisal.AppraisalKpis {
		if ak.EmployeeID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "employee_id field is required"})
			return
		}

		if ak.KpiID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "kpi_id field is required "})
			return
		}

		if ak.Status == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "status field is required "})
			return
		}
	}

	// check appraisal type exists
	err = checkAppraisalType(r.Db, appraisal.AppraisalTypeStr)
	if err != nil {
		log.Error("invalid appraisal type")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid appraisal type"})
		return
	}

	// checking appraisal flow id exists in db
	var appraisalFlow models.AppraisalFlow
	err = r.Db.Model(&models.AppraisalFlow{}).First(&appraisalFlow, appraisal.AppraisalFlowID).Error
	if err != nil {
		log.Error("invalid appraisal flow id")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid appraisal flow id"})
		return
	}

	dbAppraisal, err := controller.CreateAppraisal(r.Db, &appraisal)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dbAppraisal)
}

func (r *AppraisalService) GetAppraisalByID(c *gin.Context) {
	log.Info("Initializing GetAppraisalByID handler function...")

	id, _ := strconv.ParseUint(c.Param("id"), 0, 64)

	var appraisal models.Appraisal
	err := controller.GetAppraisalByID(r.Db, &appraisal, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("appraisal record not found against the given id")
			c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
			return
		}

		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appraisal)
}

func (r *AppraisalService) GetAllAppraisals(c *gin.Context) {
	log.Info("Initializing GetAllAppraisal handler function...")
	var appraisals []models.Appraisal
	db := r.Db.Model(&models.Appraisal{})

	appraisalName := c.Query("appraisal_name")
	supervisorID := c.Query("supervisor_id")

	if appraisalName != "" {
		db = db.Where("appraisal_name LIKE ?", "%"+appraisalName+"%")
	}

	if supervisorID != "" {
		db = db.Where("supervisor_id = ?", supervisorID)
	}

	err := controller.GetAllAppraisals(db, &appraisals)

	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appraisals)
}

func (r *AppraisalService) UpdateAppraisal(c *gin.Context) {
	log.Info("Initializing UpdateAppraisal handler function...")

	id, _ := strconv.ParseUint(c.Param("id"), 0, 64)

	var appraisal models.Appraisal

	err := c.ShouldBindJSON(&appraisal)
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

	// Validate each FlowStep struct
	for _, ak := range appraisal.AppraisalKpis {
		if ak.EmployeeID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "employee_id field is required"})
			return
		}

		if ak.KpiID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "kpi_id field is required "})
			return
		}

		if ak.Status == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "status field is required "})
			return
		}
	}

	// check appraisal type exists
	err = checkAppraisalType(r.Db, appraisal.AppraisalTypeStr)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid appraisal type"})
		return
	}
	appraisal.ID = id

	// checking appraisal flow id exists in db
	var appraisalFlow models.AppraisalFlow
	err = r.Db.Model(&models.AppraisalFlow{}).First(&appraisalFlow, appraisal.AppraisalFlowID).Error
	if err != nil {
		log.Error("invalid appraisal flow id")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid appraisal flow id"})
		return
	}

	// callling controller update function
	dbAppraisal, err := controller.UpdateAppraisal(r.Db, &appraisal)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dbAppraisal)
}

func (r *AppraisalService) DeleteAppraisal(c *gin.Context) {
	log.Info("Initializing DeleteAppraisal handler function...")

	var appraisal models.Appraisal
	id, _ := strconv.ParseUint(c.Param("id"), 0, 64)
	appraisal.ID = id

	err := controller.GetAppraisalByID(r.Db, &appraisal, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("appraisal record not found against the given id")
			c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
			return
		}

		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = controller.DeleteAppraisal(r.Db, &appraisal, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func checkAppraisalType(db *gorm.DB, appraisal_type string) error {
	log.Info("Checking Appraisal type")
	var appraisalTypeModel models.AppraisalType
	err := db.Model(&models.AppraisalType{}).Where("appraisal_type = ?", appraisal_type).First(&appraisalTypeModel).Error
	if err != nil {
		return err
	}

	return nil
}
