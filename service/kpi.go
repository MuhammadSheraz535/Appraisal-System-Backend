// Service/kpi.go

package service

import (
	"net/http"
	"strconv"

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
	SINGLE_KPI_TYPE = "Single"
	MULTI_KPI_TYPE  = "Multi"
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
	// TODO: Delete this table and get KPI types from /kpi_types endpoint
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
			newKpiType.BasicKpiType = SINGLE_KPI_TYPE
		} else if k == 2 || k == 3 {
			newKpiType.BasicKpiType = MULTI_KPI_TYPE
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

	// Make assign type ID as 0 = Role, 1 = Team and 2 = Individual
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

	kpi.ID = 0

	kpiType, err := checkKpiType(s.Db, kpi.KpiType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid KPI type"})
		return
	}

	_, err = checkAssignType(s.Db, kpi.AssignType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assign type"})
		return
	}

	switch kpiType.BasicKpiType {
	case SINGLE_KPI_TYPE:
		if kpi.Statement == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "statement is nil"})
			return
		}

		kpi, err = controller.CreateKPI(s.Db, kpi)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, kpi)

	case MULTI_KPI_TYPE:
		var multiKpi models.MultiKpi
		err = c.ShouldBindBodyWith(&multiKpi, binding.JSON)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		tx := s.Db.Begin()
		kpi.Statement = "" // since it's multi-statement
		kpi, err = controller.CreateKPI(tx, kpi)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		multiKpi.ID = kpi.ID

		for i, statement := range multiKpi.Statements {
			statement.KpiID = kpi.ID
			err = tx.Create(&statement).Error
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			multiKpi.Statements[i].ID = statement.ID
			multiKpi.Statements[i].KpiID = kpi.ID
		}

		err = tx.Commit().Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, multiKpi)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid KPI type"})
		return
	}
}

func (s *KPIService) UpdateKPI(c *gin.Context) {
	kpiID := c.Param("id")
	var kpi models.Kpi
	var err error

	if err := c.ShouldBindBodyWith(&kpi, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.Atoi(kpiID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	kpi.ID = uint(id)

	kpiType, err := checkKpiType(s.Db, kpi.KpiType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid KPI type"})
		return
	}

	_, err = checkAssignType(s.Db, kpi.AssignType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assign type"})
		return
	}

	switch kpiType.BasicKpiType {
	case SINGLE_KPI_TYPE:
		if kpi.Statement == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Statement is nil"})
			return
		}

		kpi, err = controller.UpdateKPI(s.Db, kpi)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, kpi)

	case MULTI_KPI_TYPE:
		var multiKpi models.MultiKpi
		err = c.ShouldBindBodyWith(&multiKpi, binding.JSON)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		tx := s.Db.Begin()
		kpi.Statement = ""
		kpi, err = controller.UpdateKPI(tx, kpi)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		multiKpi.ID = kpi.ID

		err = tx.Where("kpi_id = ?", kpi.ID).Delete(&models.MultiStatementKpiData{}).Error
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for i, statement := range multiKpi.Statements {
			statement.KpiID = kpi.ID
			err = tx.Create(&statement).Error
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			multiKpi.Statements[i].ID = statement.ID
			multiKpi.Statements[i].KpiID = kpi.ID
		}

		err = tx.Commit().Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, multiKpi)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid KPI type"})
		return
	}
}

func (s *KPIService) GetKPIByID(c *gin.Context) {
	kpiID := c.Param("id")
	id, err := strconv.Atoi(kpiID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	kpi, err := controller.GetKPIByID(s.Db, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	kpiType, err := checkKpiType(s.Db, kpi.KpiType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	switch kpiType.BasicKpiType {
	case SINGLE_KPI_TYPE:
		kpi_data := models.Kpi{
			ID:            kpi.ID,
			KpiName:       kpi.KpiName,
			AssignType:    kpi.AssignType,
			KpiType:       kpi.KpiType,
			ApplicableFor: kpi.ApplicableFor,
			Statement:     kpi.Statement,
		}

		err := s.Db.Find(&kpi_data).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch KPI"})
			return
		}
		c.JSON(http.StatusOK, kpi_data)

	case MULTI_KPI_TYPE:
		kpi_data := models.MultiKpi{
			ID:            kpi.ID,
			KpiName:       kpi.KpiName,
			AssignType:    kpi.AssignType,
			KpiType:       kpi.KpiType,
			ApplicableFor: kpi.ApplicableFor,
		}

		var multistatementKpi []models.MultiStatementKpiData
		err := s.Db.Where("kpi_id = ?", kpi.ID).Find(&multistatementKpi).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch KPI"})
			return
		}
		kpi_data.Statements = append(kpi_data.Statements, multistatementKpi...)

		c.JSON(http.StatusOK, kpi_data)
	}
}

func (s *KPIService) GetAllKPI(c *gin.Context) {
	var kpis []models.Kpi
	db := s.Db

	kpiName := c.Query("kpi_name")
	assignType := c.Query("assign_type")
	kpiType := c.Query("kpi_type")

	if kpiName != "" && assignType != "" && kpiType != "" {
		db = db.Table("kpis").Where("kpis.kpi_name LIKE ? AND kpis.assign_type = ? AND kpis.kpi_type = ?", "%"+kpiName+"%", assignType, kpiType)
	} else if kpiName != "" && assignType != "" {
		db = db.Table("kpis").Where("kpis.kpi_name LIKE ? AND kpis.assign_type = ?", "%"+kpiName+"%", assignType)
	} else if kpiName != "" && kpiType != "" {
		db = db.Table("kpis").Where("kpis.kpi_name LIKE ? AND kpis.kpi_type = ?", "%"+kpiName+"%", kpiType)
	} else if assignType != "" && kpiType != "" {
		db = db.Table("kpis").Where("kpis.assign_type = ? AND kpis.kpi_type = ?", assignType, kpiType)
	} else if kpiName != "" {
		db = db.Table("kpis").Where("kpis.kpi_name LIKE ?", "%"+kpiName+"%")
	} else if assignType != "" {
		db = db.Table("kpis").Where("kpis.assign_type = ?", assignType)
	} else if kpiType != "" {
		db = db.Table("kpis").Where("kpis.kpi_type = ?", kpiType)
	}

	if err := controller.GetAllKPI(db, &kpis); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch KPIs"})
		return
	}

	var allKpis []interface{}

	for _, k := range kpis {
		kpiType, err := checkKpiType(s.Db, k.KpiType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		switch kpiType.BasicKpiType {
		case SINGLE_KPI_TYPE:
			kpi_data := models.Kpi{
				ID:            k.ID,
				KpiName:       k.KpiName,
				AssignType:    k.AssignType,
				KpiType:       k.KpiType,
				ApplicableFor: k.ApplicableFor,
				Statement:     k.Statement,
			}
			err := s.Db.Find(&kpis).Error
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch KPI"})
				return
			}
			allKpis = append(allKpis, kpi_data)

		case MULTI_KPI_TYPE:
			kpi_data := models.MultiKpi{
				ID:            k.ID,
				KpiName:       k.KpiName,
				AssignType:    k.AssignType,
				KpiType:       k.KpiType,
				ApplicableFor: k.ApplicableFor,
			}

			var multistatementKpi []models.MultiStatementKpiData
			err := s.Db.Where("kpi_id = ?", k.ID).Find(&multistatementKpi).Error
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch KPI"})
				return
			}

			kpi_data.Statements = append(kpi_data.Statements, multistatementKpi...)

			allKpis = append(allKpis, kpi_data)
		}
	}

	c.JSON(http.StatusOK, allKpis)
}

// DeleteKPI deletes a KPI with the given ID
func (s *KPIService) DeleteKPI(c *gin.Context) {
	id := c.Param("id")

	var count int64
	var kpi models.Kpi
	if err := s.Db.First(&kpi, id).Count(&count).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "KPI not found"})
		return
	}

	tx := s.Db.Begin()

	if err := controller.DeleteKPI(tx, id); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	kpiType, err := checkKpiType(tx, kpi.KpiType)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if kpiType.BasicKpiType == MULTI_KPI_TYPE {
		if err := tx.Where("kpi_id = ?", kpi.ID).Delete(&models.MultiStatementKpiData{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete KPI statements"})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.Status(http.StatusNoContent)
}

func checkKpiType(db *gorm.DB, kpiType string) (models.KpiType, error) {
	var kpiTypeModel models.KpiType
	err := db.Where("kpi_type = ?", kpiType).First(&kpiTypeModel).Error
	if err != nil {
		return kpiTypeModel, err
	}
	return kpiTypeModel, nil
}

func checkAssignType(db *gorm.DB, assignType uint64) (models.AssignType, error) {
	var assignTypeModel models.AssignType
	err := db.Where("assign_type_id = ?", assignType).First(&assignTypeModel).Error
	if err != nil {
		return assignTypeModel, err
	}
	return assignTypeModel, nil
}
