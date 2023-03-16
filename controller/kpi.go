package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/mrehanabbasi/appraisal-system-backend/database"
	"github.com/mrehanabbasi/appraisal-system-backend/models"
)

type KPIController struct {
	Db *gorm.DB
}

func NewKPIController() *KPIController {
	db := database.DB
	db.AutoMigrate(&models.KPI{}, &models.MeasuredData{}, &models.QuestionaireData{})
	return &KPIController{Db: db}
}

func CreateKPI(db *gorm.DB, KPI models.KPI) (models.KPI, error) {
	// Serialize the KPI struct to JSON
	data, err := json.Marshal(KPI)
	if err != nil {
		return KPI, err
	}

	// Deserialize the JSON to a new KPI struct with StringSlice fields
	var newKPI models.KPI
	err = json.Unmarshal(data, &newKPI)
	if err != nil {
		return KPI, err
	}

	// Create the new KPI in the database
	err = db.Table("kpis").Create(&newKPI).Error
	if err != nil {
		return KPI, err
	}

	return newKPI, nil
}

// CreateKPI creates a new KPI
func (r *KPIController) CreateKPI(c *gin.Context) {
	kpi := models.KPI{}

	// Bind request body to KPI model
	if err := c.BindJSON(&kpi); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Convert RolesApplicable to JSON strings

	rolesApplicableJSON, err := json.Marshal(kpi.RolesApplicable)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	kpi.RolesApplicable = nil

	// Create new KPI
	if err := r.Db.Create(&kpi).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Set RolesApplicable fields to original string slices

	kpi.RolesApplicable = []string{}
	if err := json.Unmarshal(rolesApplicableJSON, &kpi.RolesApplicable); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, kpi)
}

// GetAllKPI retrieves all KPIs from the database
func GetAllKPI(db *gorm.DB) ([]models.KPI, error) {
	kpis := []models.KPI{}

	// Retrieve all KPIs from the database
	if err := db.Table("kpis").Find(&kpis).Error; err != nil {
		return nil, err
	}

	// Convert RolesApplicable to JSON strings for serialization
	for i := range kpis {

		rolesApplicableJSON, err := json.Marshal(kpis[i].RolesApplicable)
		if err != nil {
			return nil, err
		}
		kpis[i].RolesApplicable = nil

		// Set RolesApplicable fields to original string slices

		kpis[i].RolesApplicable = []string{}
		if err := json.Unmarshal(rolesApplicableJSON, &kpis[i].RolesApplicable); err != nil {
			return nil, err
		}
	}

	return kpis, nil
}

// GetAllKPI retrieves all KPIs from the database
func (r *KPIController) GetAllKPI(c *gin.Context) {
	kpis := []models.KPI{}

	// Retrieve query parameters
	kpiName := c.Query("KPIName")
	assignType := c.Query("AssignType")
	rolesApplicable := c.Query("RolesApplicable")

	// Build query based on query parameters
	query := r.Db
	if kpiName != "" {
		query = query.Where("kpi_name LIKE ?", fmt.Sprintf("%%%s%%", kpiName))
	}
	if assignType != "" {
		query = query.Where("assign_type LIKE ?", fmt.Sprintf("%%%s%%", assignType))
	}
	if rolesApplicable != "" {
		query = query.Where("roles_applicable LIKE ?", fmt.Sprintf("%%%s%%", rolesApplicable))
	}

	// Retrieve all KPIs from the database that match the query
	if err := query.Find(&kpis).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Convert RolesApplicable to JSON strings for serialization
	for i := range kpis {

		rolesApplicableJSON, err := json.Marshal(kpis[i].RolesApplicable)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		kpis[i].RolesApplicable = nil

		// Set RolesApplicable fields to original string slices

		kpis[i].RolesApplicable = []string{}
		if err := json.Unmarshal(rolesApplicableJSON, &kpis[i].RolesApplicable); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	c.JSON(http.StatusOK, kpis)
}

// GetKPIByID retrieves a KPI by its ID from the database
func GetKPIByID(db *gorm.DB, id string) (models.KPI, error) {
	var kpi models.KPI
	if err := db.Table("kpis").Where("id = ?", id).First(&kpi).Error; err != nil {
		return kpi, err
	}

	// Convert RolesApplicable to JSON strings

	rolesApplicableJSON, err := json.Marshal(kpi.RolesApplicable)
	if err != nil {
		return kpi, err
	}
	kpi.RolesApplicable = nil

	// Set RolesApplicable fields to original string slices

	kpi.RolesApplicable = []string{}
	if err := json.Unmarshal(rolesApplicableJSON, &kpi.RolesApplicable); err != nil {
		return kpi, err
	}

	return kpi, nil
}

// GetKPIByID retrieves a KPI by its ID
func (r *KPIController) GetKPIByID(c *gin.Context) {
	id := c.Param("id")

	var kpi models.KPI
	if err := r.Db.Where("id = ?", id).First(&kpi).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	// Convert RolesApplicable to JSON strings

	rolesApplicableJSON, err := json.Marshal(kpi.RolesApplicable)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	kpi.RolesApplicable = nil

	// Set RolesApplicable fields to original string slices

	kpi.RolesApplicable = []string{}
	if err := json.Unmarshal(rolesApplicableJSON, &kpi.RolesApplicable); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, kpi)
}

func UpdateKPI(db *gorm.DB, kpiID uint, updatedKPI models.KPI) (models.KPI, error) {
	// Serialize the updated KPI struct to JSON
	data, err := json.Marshal(updatedKPI)
	if err != nil {
		return updatedKPI, err
	}

	// Deserialize the JSON to a new KPI struct with StringSlice fields
	var newKPI models.KPI
	err = json.Unmarshal(data, &newKPI)
	if err != nil {
		return updatedKPI, err
	}

	// Update the KPI in the database
	err = db.Table("kpis").Where("id = ?", kpiID).Updates(newKPI).Error
	if err != nil {
		return updatedKPI, err
	}

	return newKPI, nil
}

// UpdateKPI updates an existing KPI
func (r *KPIController) UpdateKPI(c *gin.Context) {
	kpiID := c.Param("id")

	// Get existing KPI from database
	var existingKPI models.KPI
	if err := r.Db.Where("id = ?", kpiID).First(&existingKPI).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Bind request body to updated KPI model
	var updatedKPI models.KPI
	if err := c.BindJSON(&updatedKPI); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	rolesApplicableJSON, err := json.Marshal(updatedKPI.RolesApplicable)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	updatedKPI.RolesApplicable = nil

	// Update existing KPI with new data
	if _, err := UpdateKPI(r.Db, existingKPI.ID, updatedKPI); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Set RolesApplicable fields to

	updatedKPI.RolesApplicable = []string{}
	if err := json.Unmarshal(rolesApplicableJSON, &updatedKPI.RolesApplicable); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, updatedKPI)
}

func DeleteKPI(db *gorm.DB, id uint) error {
	// Delete KPI with matching ID
	err := db.Table("kpis").Where("id = ?", id).Delete(&models.KPI{}).Error
	if err != nil {
		return err
	}
	return nil
}

// DeleteKPI deletes a KPI with the given ID
func (r *KPIController) DeleteKPI(c *gin.Context) {
	// Get ID from path parameter
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Delete KPI from database
	err = DeleteKPI(r.Db, uint(id))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
