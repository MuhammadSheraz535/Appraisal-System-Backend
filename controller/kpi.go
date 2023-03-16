package controller

import (
	"net/http"

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
	db.AutoMigrate(&models.KPI{})
	return &KPIController{Db: db}
}

func GetKPIs(db *gorm.DB, KPI *[]models.KPI) (err error) {
	err = db.Table("kpis").Find(&KPI).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *KPIController) GetKPIs(c *gin.Context) {
	name := c.Query("name")
	assignType := c.Query("assign_type")
	kpiType := c.Query("type")
	role := c.Query("role")

	// Get KPIs based on query parameters
	kpis := []models.KPI{}
	query := r.Db
	if name != "" {
		query = query.Where("kpi_name = ?", name)
	}
	if assignType != "" {
		query = query.Where("assign_type = ?", assignType)
	}
	if kpiType != "" {
		query = query.Where("type = ?", kpiType)
	}
	if role != "" {
		query = query.Where("roles_applicable LIKE ?", "%"+role+"%")
	}
	if err := query.Find(&kpis).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, kpis)
}

func GetKPIByID(db *gorm.DB, KPI *models.KPI, id int) (err error) {
	err = db.Table("kpis").Where("id = ?", id).First(&KPI).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *KPIController) GetKPIByID(c *gin.Context) {
	id := c.Param("id")

	// Get KPI by ID
	kpi := models.KPI{}
	if err := r.Db.Where("id = ?", id).First(&kpi).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, kpi)
}

func CreateKPI(db *gorm.DB, KPI models.KPI) (models.KPI, error) {
	err := db.Table("kpis").Create(&KPI).Error
	if err != nil {
		return KPI, err
	}
	return KPI, nil
}

// CreateKPI creates a new KPI
func (r *KPIController) CreateKPI(c *gin.Context) {
	kpi := models.KPI{}

	// Bind request body to KPI model
	if err := c.BindJSON(&kpi); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Create KPI
	if err := r.Db.Create(&kpi).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, kpi)
}

func UpdateKPI(db *gorm.DB, KPI *models.KPI) error {
	err := db.Table("kpis").Updates(KPI).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateKPI updates an existing KPI
func (r *KPIController) UpdateKPI(c *gin.Context) {
	id := c.Param("id")

	kpi := models.KPI{}

	// Get KPI by ID
	if err := r.Db.Where("id = ?", id).First(&kpi).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Bind request body to KPI model and update fields
	if err := c.BindJSON(&kpi); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Update KPI
	if err := r.Db.Save(&kpi).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, kpi)
}

func DeleteKPI(db *gorm.DB, KPI *models.KPI, id int) error {
	err := db.Table("kpis").Where("id = ?", id).Delete(&KPI).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *KPIController) DeleteKPI(c *gin.Context) {
	id := c.Param("id")
	kpi := models.KPI{}
	// Get KPI by ID
	if err := r.Db.Where("id = ?", id).First(&kpi).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	// Delete KPI
	if err := r.Db.Delete(&kpi).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusNoContent)
}
