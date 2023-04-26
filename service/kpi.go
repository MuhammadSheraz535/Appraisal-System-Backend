// Service/kpi.go

package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/mrehanabbasi/appraisal-system-backend/controller"
	"github.com/mrehanabbasi/appraisal-system-backend/database"
	log "github.com/mrehanabbasi/appraisal-system-backend/logger"
	"github.com/mrehanabbasi/appraisal-system-backend/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	// TODO: Delete this table population and get KPI types from /kpi_types endpoint
	kpiTypes := []string{
		FEEDBACK_KPI_TYPE,
		OBSERVATORY_KPI_TYPE,
		MEASURED_KPI_TYPE,
		QUESTIONNAIRE_KPI_TYPE,
	}

	kpiTypesSlice := make([]models.KpiType, len(kpiTypes))

	for k, v := range kpiTypes {
		newKpiType := models.KpiType{
			KpiType: v,
		}
		if k == 0 || k == 1 { // Feedback and Observatory will be 'Single'
			newKpiType.BasicKpiType = SINGLE_KPI_TYPE
		} else if k == 2 || k == 3 { // Measured and Questionnaire will be 'Multi'
			newKpiType.BasicKpiType = MULTI_KPI_TYPE
		}

		kpiTypesSlice[k] = newKpiType
	}

	err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&kpiTypesSlice).Error
	if err != nil {
		log.Error(err.Error())
		return err
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
	assignTypesSlice := make([]models.AssignType, len(assignTypes))
	for i, a := range assignTypes {
		newAssignType := models.AssignType{
			AssignTypeId: uint64(i),
			AssignType:   a,
		}
		assignTypesSlice[i] = newAssignType
	}

	err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&assignTypesSlice).Error
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

func (s *KPIService) CreateKPI(c *gin.Context) {
	log.Info("Initializing CreateKPI handler function...")

	var kpi models.Kpi
	var err error

	if err := c.ShouldBindBodyWith(&kpi, binding.JSON); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	kpi.ID = 0

	kpiType, err := checkKpiType(s.Db, kpi.KpiTypeID)
	if err != nil {
		log.Error("invalid kpi type")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid KPI type"})
		return
	}

	_, err = checkAssignType(s.Db, kpi.AssignTypeID)
	if err != nil {
		log.Error("invalid assign type")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assign type"})
		return
	}

	switch kpiType.BasicKpiType {
	case SINGLE_KPI_TYPE:
		log.Info("Handling single KPI type request...")

		if kpi.Statement == "" {
			log.Error("statement is nil in the request")
			c.JSON(http.StatusBadRequest, gin.H{"error": "statement is nil"})
			return
		}

		kpi, err = controller.CreateKPI(s.Db, kpi)
		if err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, kpi)

	case MULTI_KPI_TYPE:
		log.Info("Handling multi KPI type request...")

		var multiKpi models.MultiKpi
		err = c.ShouldBindBodyWith(&multiKpi, binding.JSON)
		if err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		tx := s.Db.Begin()
		kpi.Statement = "" // since it's multi-statement
		kpi, err = controller.CreateKPI(tx, kpi)
		if err != nil {
			tx.Rollback()
			log.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		multiKpi.ID = kpi.ID

		for i, statement := range multiKpi.Statements {
			statement.KpiID = kpi.ID
			err = tx.Create(&statement).Error
			if err != nil {
				tx.Rollback()
				log.Error(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			multiKpi.Statements[i].ID = statement.ID
			multiKpi.Statements[i].KpiID = kpi.ID
		}

		err = tx.Commit().Error
		if err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, multiKpi)
	default:
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid KPI type"})
		return
	}
}

func (s *KPIService) UpdateKPI(c *gin.Context) {
	log.Info("Initializing UpdateKPI handler function...")

	kpiID := c.Param("id")
	var kpi models.Kpi
	var err error

	if err := c.ShouldBindBodyWith(&kpi, binding.JSON); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.Atoi(kpiID)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	kpi.ID = uint64(id)

	kpiType, err := checkKpiType(s.Db, kpi.KpiTypeID)
	if err != nil {
		log.Error("invalid kpi type")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid KPI type"})
		return
	}

	_, err = checkAssignType(s.Db, kpi.AssignTypeID)
	if err != nil {
		log.Error("invalid assign type")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assign type"})
		return
	}

	switch kpiType.BasicKpiType {
	case SINGLE_KPI_TYPE:
		log.Info("Handling single KPI type request...")

		if kpi.Statement == "" {
			log.Error("statement is nil in the request")
			c.JSON(http.StatusBadRequest, gin.H{"error": "statement is nil"})
			return
		}

		kpi, err = controller.UpdateKPI(s.Db, kpi)
		if err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, kpi)

	case MULTI_KPI_TYPE:
		log.Info("Handling multi KPI type request...")

		var multiKpi models.MultiKpi
		err = c.ShouldBindBodyWith(&multiKpi, binding.JSON)
		if err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		tx := s.Db.Begin()
		kpi.Statement = ""
		kpi, err = controller.UpdateKPI(tx, kpi)
		if err != nil {
			tx.Rollback()
			log.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		multiKpi.ID = kpi.ID

		err = tx.Where("kpi_id = ?", kpi.ID).Delete(&models.MultiStatementKpiData{}).Error
		if err != nil {
			tx.Rollback()
			log.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for i, statement := range multiKpi.Statements {
			statement.KpiID = kpi.ID
			err = tx.Create(&statement).Error
			if err != nil {
				tx.Rollback()
				log.Error(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			multiKpi.Statements[i].ID = statement.ID
			multiKpi.Statements[i].KpiID = kpi.ID
		}

		err = tx.Commit().Error
		if err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, multiKpi)
	default:
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid KPI type"})
		return
	}
}

func (s *KPIService) GetKPIByID(c *gin.Context) {
	log.Info("Initializing GetKPIByID handler function...")

	kpiID := c.Param("id")
	id, err := strconv.Atoi(kpiID)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	kpi, err := controller.GetKPIByID(s.Db, uint(id))
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	kpiType, err := checkKpiType(s.Db, kpi.KpiTypeID)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	switch kpiType.BasicKpiType {
	case SINGLE_KPI_TYPE:
		log.Info("Handling single KPI type request...")

		kpi_data := models.Kpi{
			CommonModel:   models.CommonModel{ID: kpi.ID},
			KpiName:       kpi.KpiName,
			AssignTypeID:  kpi.AssignTypeID,
			KpiTypeID:     kpi.KpiTypeID,
			ApplicableFor: kpi.ApplicableFor,
			Statement:     kpi.Statement,
		}

		err := s.Db.Find(&kpi_data).Error
		if err != nil {
			log.Error("failed to fetch single type kpi")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch KPI"})
			return
		}
		c.JSON(http.StatusOK, kpi_data)

	case MULTI_KPI_TYPE:
		log.Info("Handling multi KPI type request...")

		kpi_data := models.MultiKpi{
			CommonModel:   models.CommonModel{ID: kpi.ID},
			KpiName:       kpi.KpiName,
			AssignType:    kpi.AssignTypeID,
			KpiType:       kpi.KpiTypeID,
			ApplicableFor: kpi.ApplicableFor,
		}

		var multistatementKpi []models.MultiStatementKpiData
		err := s.Db.Where("kpi_id = ?", kpi.ID).Find(&multistatementKpi).Error
		if err != nil {
			log.Error("failed to fetch multi type kpi")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch KPI"})
			return
		}
		kpi_data.Statements = append(kpi_data.Statements, multistatementKpi...)

		c.JSON(http.StatusOK, kpi_data)
	}
}

func (s *KPIService) GetAllKPI(c *gin.Context) {
	log.Info("Initializing GetAllKPI handler function...")

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
		log.Error("failed to fetch kpis")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch KPIs"})
		return
	}

	var allKpis []interface{}

	for _, k := range kpis {
		kpiType, err := checkKpiType(s.Db, k.KpiTypeID)
		if err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		switch kpiType.BasicKpiType {
		case SINGLE_KPI_TYPE:
			kpi_data := models.Kpi{
				CommonModel:   models.CommonModel{ID: k.ID},
				KpiName:       k.KpiName,
				AssignTypeID:  k.AssignTypeID,
				KpiTypeID:     k.KpiTypeID,
				ApplicableFor: k.ApplicableFor,
				Statement:     k.Statement,
			}
			err := s.Db.Find(&kpis).Error
			if err != nil {
				log.Error("failed to fetch single type kpi")
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch KPI"})
				return
			}
			allKpis = append(allKpis, kpi_data)

		case MULTI_KPI_TYPE:
			kpi_data := models.MultiKpi{
				CommonModel:   models.CommonModel{ID: k.ID},
				KpiName:       k.KpiName,
				AssignType:    k.AssignTypeID,
				KpiType:       k.KpiTypeID,
				ApplicableFor: k.ApplicableFor,
			}

			var multistatementKpi []models.MultiStatementKpiData
			err := s.Db.Where("kpi_id = ?", k.ID).Find(&multistatementKpi).Error
			if err != nil {
				log.Error("failed to fetch multi type kpi")
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
	log.Info("Initializing DeleteKPI handler function...")

	id := c.Param("id")

	var count int64
	var kpi models.Kpi
	if err := s.Db.First(&kpi, id).Count(&count).Error; err != nil {
		log.Error("kpi not found with the given id")
		c.JSON(http.StatusNotFound, gin.H{"error": "KPI not found"})
		return
	}

	tx := s.Db.Begin()

	if err := controller.DeleteKPI(tx, id); err != nil {
		tx.Rollback()
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	kpiType, err := checkKpiType(tx, kpi.KpiTypeID)
	if err != nil {
		tx.Rollback()
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if kpiType.BasicKpiType == MULTI_KPI_TYPE {
		if err = tx.Where("kpi_id = ?", kpi.ID).Delete(&models.MultiStatementKpiData{}).Error; err != nil {
			tx.Rollback()
			log.Error("failed to delete kpi statements")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete KPI statements"})
			return
		}
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Error("failed to commit transactions")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.Status(http.StatusNoContent)
}

func checkKpiType(db *gorm.DB, kpiType string) (models.KpiType, error) {
	log.Info("Checking KPI type")
	var kpiTypeModel models.KpiType
	err := db.Where("kpi_type = ?", kpiType).First(&kpiTypeModel).Error
	if err != nil {
		return kpiTypeModel, err
	}
	return kpiTypeModel, nil
}

func checkAssignType(db *gorm.DB, assignType uint64) (models.AssignType, error) {
	log.Info("Checking assign type")
	var assignTypeModel models.AssignType
	err := db.Where("assign_type_id = ?", assignType).First(&assignTypeModel).Error
	if err != nil {
		return assignTypeModel, err
	}
	return assignTypeModel, nil
}
