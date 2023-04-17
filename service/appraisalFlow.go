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
	err := db.AutoMigrate(&models.AppraisalFlow{}, &models.FlowStep{})
	if err != nil {
		panic(err)
	}
	return &ApprasialFlowService{Db: db}
}

func (r *ApprasialFlowService) CreateAppraisalFlow(c *gin.Context) {
	var appraisalFlow models.AppraisalFlow
	err := c.ShouldBindJSON(&appraisalFlow)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appraisalFlow, err = controller.CreateAppraisalFlow(r.Db, appraisalFlow)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, appraisalFlow)
}

func (r *ApprasialFlowService) GetAppraisalFlowByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var appraisalFlow models.AppraisalFlow
	err := controller.GetAppraisalFlowByID(r.Db, &appraisalFlow, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "record not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, appraisalFlow)
}

func (r *ApprasialFlowService) GetAllApprasialFlow(c *gin.Context) {
	var appraisalFlow []models.AppraisalFlow

	flowName := c.Query("flow_name")
	isActive := c.Query("is_active")
	teamId := c.Query("team_id")

	err := controller.GetAllApprasialFlow(flowName, isActive, teamId, r.Db, &appraisalFlow)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appraisalFlow)
}

func (r *ApprasialFlowService) UpdateAppraisalFlow(c *gin.Context) {
	var appraisalFlow models.AppraisalFlow
	id, _ := strconv.Atoi(c.Param("id"))
	err := controller.GetAppraisalFlowByID(r.Db, &appraisalFlow, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	err = c.ShouldBindJSON(&appraisalFlow)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = controller.UpdateAppraisalFlow(r.Db, &appraisalFlow, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, appraisalFlow)
}

func (r *ApprasialFlowService) DeleteApprasialFlow(c *gin.Context) {
	var appraisalFlow models.AppraisalFlow
	id, _ := strconv.Atoi(c.Param("id"))
	appraisalFlow.ID = uint64(id)
	err := controller.GetAppraisalFlowByID(r.Db, &appraisalFlow, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "record not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = controller.DeleteApprasialFlow(r.Db, &appraisalFlow, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
