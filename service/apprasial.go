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
	"gorm.io/gorm/clause"
)

const (
	MID_YEAR_APPRASIAL = "MidYear"
	ANNUAL_APPRAISAL   = "Annual"
)

type ApprasialService struct {
	Db *gorm.DB
}

func NewApprasialService() *ApprasialService {
	db := database.DB
	err := db.AutoMigrate(&models.Apprasial{}, models.AppraisalType{}, models.AppraisalKpis{})
	if err != nil {
		panic(err)
	}
	// Populate appraisal_types table
	err = populateAppraisalTypeTable(db)
	if err != nil {
		panic(err)
	}
	return &ApprasialService{Db: db}
}

func populateAppraisalTypeTable(db *gorm.DB) error {
	appraisalTypes := []string{
		MID_YEAR_APPRASIAL,
		ANNUAL_APPRAISAL,
	}

	appraisalTypesSlice := make([]models.AppraisalType, len(appraisalTypes))

	for k, v := range appraisalTypes {
		newAppraisalType := models.AppraisalType{
			AppraisalType: v,
		}
		if k == 0 {
			newAppraisalType.AppraisalType = MID_YEAR_APPRASIAL
		} else if k == 2 {
			newAppraisalType.AppraisalType = ANNUAL_APPRAISAL
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

func (r *ApprasialService) CreateAppraisal(c *gin.Context) {
	log.Info("Initializing CreateAppraisal handler function...")

	var appraisal models.Apprasial
	err := c.ShouldBindJSON(&appraisal)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//check appraisal type exists
	err = checkAppraisalType(r.Db, appraisal.AppraisalTypeStr)
	if err != nil {
		log.Error("invalid Appraisal type")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid Apprasial type"})
		return
	}

	//checking appraisal flow id exists in db
	var apprasialflow models.AppraisalFlow

	err = r.Db.First(&apprasialflow, appraisal.AppraisalFlowID).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Appraisal Flow ID"})
		return
	}

	appraisal, err = controller.CreateAppraisal(r.Db, appraisal)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, appraisal)
}


func (r *ApprasialService) GetAppraisalByID(c *gin.Context) {
	log.Info("Initializing GetAppraisalByID handler function...")

	id, _ := strconv.Atoi(c.Param("id"))
	var appraisal models.Apprasial
	err := controller.GetAppraisalByID(r.Db, &appraisal, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("appraisal record not found against the given id")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "record not found"})
			return
		}

		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, appraisal)
}

func (r *ApprasialService) GetAllApprasial(c *gin.Context) {
	log.Info("Initializing GetAllAppraisal handler function...")
	var appraisal []models.Apprasial
	db := r.Db.Model(&models.Apprasial{})

	apprasialName := c.Query("apprasial_name")
	supervisorID := c.Query("supervisor_id")

	if apprasialName != "" {
		db = db.Where("apprasial_name LIKE ?", "%"+apprasialName+"%")
	}

	if supervisorID != "" {
		db = db.Where("supervisor_id = ?", supervisorID)
	}

	err := controller.GetAllApprasial(db, &appraisal)

	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appraisal)
}

func (r *ApprasialService) UpdateAppraisal(c *gin.Context) {
	log.Info("Initializing UpdateAppraisal handler function...")

	var appraisal models.Apprasial
	id, _ := strconv.Atoi(c.Param("id"))

	err := controller.GetAppraisalByID(r.Db, &appraisal, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("appraisal record not found against the given id")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "record not found"})
			return
		}

		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = c.ShouldBindJSON(&appraisal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Query apprasialflow struct using apprasialflowid value and update  apprasialflow id accordingly
	var apprasialflow models.AppraisalFlow
	err = r.Db.First(&apprasialflow, appraisal.AppraisalFlowID).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Appraisal Flow ID not exist"})
		return
	}
	appraisal.AppraisalFlow = apprasialflow

	appraisal, err = controller.UpdateAppraisal(r.Db, &appraisal, id)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, appraisal)
}

func (r *ApprasialService) DeleteApprasial(c *gin.Context) {
	log.Info("Initializing DeleteAppraisal handler function...")

	var appraisal models.Apprasial
	id, _ := strconv.Atoi(c.Param("id"))
	appraisal.ID = uint64(id)

	err := controller.GetAppraisalByID(r.Db, &appraisal, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("appraisal record not found against the given id")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "record not found"})
			return
		}

		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = controller.DeleteApprasial(r.Db, &appraisal, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func checkAppraisalType(db *gorm.DB, appraisal_type string) ( error) {
	log.Info("Checking Appraisal type")
	var appraisalTypeModel models.AppraisalType
	err := db.Where("appraisal_type = ?", appraisal_type).First(&appraisalTypeModel).Error
	if err != nil {
		return  err
	}
	return nil
}