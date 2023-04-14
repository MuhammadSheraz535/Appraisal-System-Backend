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
