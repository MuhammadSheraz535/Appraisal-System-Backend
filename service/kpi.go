// Service/kpi.go

package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/mrehanabbasi/appraisal-system-backend/controller"
	"github.com/mrehanabbasi/appraisal-system-backend/database"
	"github.com/mrehanabbasi/appraisal-system-backend/models"
	"gorm.io/gorm"
)

const (
	FEEDBACK_KPI_TYPE      = "Feedback"
	OBSERVATORY_KPI_TYPE   = "Observatory"
	MEASURED_KPI_TYPE      = "Measured"
	QUESTIONNAIRE_KPI_TYPE = "Questionnaire"
)
const (
	Single_KPI_TYPE   = "Single"
	Multiple_KPI_TYPE = "Multiple"
)

const (
	ASSIGN_TYPE_ROLE       = "Role"
	ASSIGN_TYPE_TEAM       = "Team"
	ASSIGN_TYPE_INDIVIDUAL = "Individual"
)

type KPIService struct {
	Db *gorm.DB
}

func NewKPIService() *KPIService {
	db := database.DB
	err := db.AutoMigrate(&models.Kpi{}, &models.KpiType{}, &models.AssignType{}, &models.MultiStatementKpiData{})
	if err != nil {
		panic(err)
	}
	// Populate assign_types table
	err = populateAssignTypeTable(db)
	if err != nil {
		panic(err)
	}

	// Populate kpi_types table
	err = populateKpiTypeTable(db)
	if err != nil {
		panic(err)
	}

	return &KPIService{Db: db}
}

func populateKpiTypeTable(db *gorm.DB) error {
	kpiTypes := []string{
		FEEDBACK_KPI_TYPE,
		OBSERVATORY_KPI_TYPE,
		MEASURED_KPI_TYPE,
		QUESTIONNAIRE_KPI_TYPE,
	}

	for k, v := range kpiTypes {
		newKpiType := models.KpiType{
			KpiType: v,
		}
		if k == 0 || k == 1 {
			newKpiType.BasicKpiType = "Single"
		} else if k == 2 || k == 3 {
			newKpiType.BasicKpiType = "Multiple"
		}
		err := db.Create(&newKpiType).Error
		if err != nil {
			return err
		}

	}
	return nil
}

func populateAssignTypeTable(db *gorm.DB) error {
	assignTypes := []string{
		ASSIGN_TYPE_ROLE,
		ASSIGN_TYPE_TEAM,
		ASSIGN_TYPE_INDIVIDUAL,
	}

	for i, a := range assignTypes {
		newAssignType := models.AssignType{
			AssignTypeId: uint64(i),
			AssignType:   a,
		}
		err := db.Create(&newAssignType).Error
		if err != nil {
			return err
		}
	}

	return nil
}
func (s *KPIService) CreateKPI(c *gin.Context) {
	var kpi models.Kpi
	var err error

	if err := c.ShouldBindBodyWith(&kpi, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if kpi.Statement == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Statement is nil"})
		return
	}

	kpi.ID = 0

	switch kpi.KpiType {
	case FEEDBACK_KPI_TYPE, OBSERVATORY_KPI_TYPE:

		kpi, err = controller.CreateKPI(s.Db, kpi)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, kpi)

	case MEASURED_KPI_TYPE, QUESTIONNAIRE_KPI_TYPE:
		var multiKpi models.MultiKpi
		err = c.ShouldBindBodyWith(&multiKpi, binding.JSON)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		tx := s.Db.Begin()
		err = tx.Create(&kpi).Error
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for _, statement := range multiKpi.Statements {
			statement.KpiID = kpi.ID
			err = tx.Create(&statement).Error
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		err = tx.Commit().Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, kpi)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid KPI type"})
		return
	}
}
