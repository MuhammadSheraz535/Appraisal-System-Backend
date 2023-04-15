package service

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mrehanabbasi/appraisal-system-backend/controller"
	"github.com/mrehanabbasi/appraisal-system-backend/database"
	"github.com/mrehanabbasi/appraisal-system-backend/models"
	"gorm.io/gorm"
)

type ApprasialFlowService struct {
	Db *gorm.DB
}

func NewApprasialFlowService() *ApprasialFlowService {
	db := database.DB
	db.AutoMigrate(&models.ApraisalFlow{}, &models.FlowStep{})
	return &ApprasialFlowService{Db: db}
}

func (r *ApprasialFlowService) CreateAppraisalFlow(c *gin.Context) {
	var appraisalflow models.ApraisalFlow
	err := c.ShouldBindJSON(&appraisalflow)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appraisalflow, err = controller.CreateAppraisalFlow(r.Db, appraisalflow)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, appraisalflow)
}

func (r *ApprasialFlowService) GetAppraisalFlowByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var apprasialflow models.ApraisalFlow
	err := controller.GetAppraisalFlowByID(r.Db, &apprasialflow, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Record not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, apprasialflow)
}

func (r *ApprasialFlowService) GetAllApprasialFlow(c *gin.Context) {
	var apprasialflow []models.ApraisalFlow
	var err error

	flowName := c.Query("flow_name")
	isActive := c.Query("is_active")
	teamId := c.Query("team_id")

	err = controller.GetAllApprasialFlow(flowName ,isActive,teamId,r.Db,&apprasialflow) 

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, apprasialflow)
}

func (r *ApprasialFlowService) UpdateAppraisalFlow(c *gin.Context) {
	var appraisalflow models.ApraisalFlow
	id, _ := strconv.Atoi(c.Param("id"))
	err := controller.GetAppraisalFlowByID(r.Db, &appraisalflow, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	err = c.ShouldBindJSON(&appraisalflow)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = controller.UpdateAppraisalFlow(r.Db, &appraisalflow)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, appraisalflow)
}
