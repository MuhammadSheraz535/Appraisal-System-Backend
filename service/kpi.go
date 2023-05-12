// Service/kpi.go

package service

import (
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
		constants.FEEDBACK_KPI_TYPE,
		constants.OBSERVATORY_KPI_TYPE,
		constants.MEASURED_KPI_TYPE,
		constants.QUESTIONNAIRE_KPI_TYPE,
	}

	kpiTypesSlice := make([]models.KpiType, len(kpiTypes))

	for k, v := range kpiTypes {
		newKpiType := models.KpiType{
			KpiType: v,
		}
		if k == 0 || k == 1 { // Feedback and Observatory will be 'Single'
			newKpiType.BasicKpiType = constants.SINGLE_KPI_TYPE
		} else if k == 2 || k == 3 { // Measured and Questionnaire will be 'Multi'
			newKpiType.BasicKpiType = constants.MULTI_KPI_TYPE
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
		constants.ASSIGN_TYPE_ROLE,
		constants.ASSIGN_TYPE_TEAM,
		constants.ASSIGN_TYPE_INDIVIDUAL,
	}

	// Make assign type ID as 0 = Role, 1 = Team and 2 = Individual
	assignTypesSlice := make([]models.AssignType, len(assignTypes))
	for i, a := range assignTypes {
		newAssignType := models.AssignType{
			AssignTypeId: uint64(i),
			AssignType:   models.AssignTypeStr(a),
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

	if err := c.ShouldBindJSON(&kpi); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// validate the kpi struct using the validator
	err = kpi.Validate()
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

	kpi.ID = 0

	kpiType, err := checkKpiType(s.Db, kpi.KpiTypeStr)
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
	case constants.SINGLE_KPI_TYPE:
		if kpi.Statement == "" {
			log.Error("statement is nil in the request")
			c.JSON(http.StatusBadRequest, gin.H{"error": "statement is nil"})
			return
		}

		kpi.Statements = nil
	case constants.MULTI_KPI_TYPE:
		if len(kpi.Statements) == 0 {
			log.Error("statements are nil in the request")
			c.JSON(http.StatusBadRequest, gin.H{"error": "statements field is nil"})
			return
		}

		kpi.Statement = ""
	}
	// Validate MultiStatementKpiData fields
	for _, mskd := range kpi.Statements {
		err = mskd.Validate()
		if err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

	}

	dbKpi, err := controller.CreateKPI(s.Db, &kpi)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dbKpi)
}

func (s *KPIService) UpdateKPI(c *gin.Context) {
	log.Info("Initializing UpdateKPI handler function...")

	kpiID := c.Param("id")
	var kpi models.Kpi
	var err error

	if err := c.ShouldBindJSON(&kpi); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// validate the kpi struct using the validator
	err = kpi.Validate()
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

	id, err := strconv.ParseUint(kpiID, 0, 64)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	kpi.ID = id

	kpiType, err := checkKpiType(s.Db, kpi.KpiTypeStr)
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
	case constants.SINGLE_KPI_TYPE:
		if kpi.Statement == "" {
			log.Error("statement is nil in the request")
			c.JSON(http.StatusBadRequest, gin.H{"error": "statement is nil"})
			return
		}

		// If the Kpi is being updated from a MultiStatementKpi to a SingleStatementKpi,
		// delete all existing MultiStatementKpiData records for the given KpiID.
		err = s.Db.Where("kpi_id = ?", kpi.ID).Delete(&models.MultiStatementKpiData{}).Error
		if err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		kpi.Statements = nil
	case constants.MULTI_KPI_TYPE:
		if len(kpi.Statements) == 0 {
			log.Error("statements are nil in the request")
			c.JSON(http.StatusBadRequest, gin.H{"error": "statements field is nil"})
			return
		}

		kpi.Statement = ""
	}
	// Validate MultiStatementKpiData fields
	for _, mskd := range kpi.Statements {
		err = mskd.Validate()
		if err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

	}

	dbKpi, err := controller.UpdateKPI(s.Db, &kpi)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dbKpi)
}

func (s *KPIService) GetKPIByID(c *gin.Context) {
	log.Info("Initializing GetKPIByID handler function...")

	kpiID := c.Param("id")
	id, err := strconv.ParseUint(kpiID, 0, 64)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	kpi, err := controller.GetKPIByID(s.Db, id)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, kpi)
}

func (s *KPIService) GetAllKPIs(c *gin.Context) {
	log.Info("Initializing GetAllKPI handler function...")

	var kpis []models.Kpi
	db := s.Db.Model(&models.Kpi{})

	kpiName := c.Query("kpi_name")
	assignType := c.Query("assign_type")
	kpiType := c.Query("kpi_type")

	if kpiName != "" {
		db = db.Where("kpi_name LIKE ?", "%"+kpiName+"%")
	}

	if assignType != "" {
		db = db.Where("assign_type_id = ?", assignType)
	}

	if kpiType != "" {
		db = db.Where("kpi_type_str = ?", kpiType)
	}

	if err := controller.GetAllKPI(db, &kpis); err != nil {
		log.Error("failed to fetch kpis")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch KPIs"})
		return
	}

	c.JSON(http.StatusOK, kpis)
}

// DeleteKPI deletes a KPI with the given ID
func (s *KPIService) DeleteKPI(c *gin.Context) {
	log.Info("Initializing DeleteKPI handler function...")

	kpiID := c.Param("id")
	id, err := strconv.ParseUint(kpiID, 0, 64)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := controller.DeleteKPI(s.Db, id); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
