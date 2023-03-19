package controller

import (
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
	db.AutoMigrate(&models.KPI{}, models.MeasuredData{}, models.QuestionaireData{})
	return &KPIController{Db: db}
}

func CreateKPI(db *gorm.DB, kpi models.KPI) (models.KPI, error) {
	if err := db.Create(&kpi).Error; err != nil {
		return kpi, err
	}
	return kpi, nil
}

func (r *KPIController) CreateKPI(c *gin.Context) {
	var kpi models.KPI
	err := c.ShouldBindJSON(&kpi)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	kpi, err = CreateKPI(r.Db, kpi)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, kpi)
}

func GetAllKPI(db *gorm.DB, queryParams map[string]string) ([]models.KPI, error) {
	var kpis []models.KPI

	query := db
	if kpiName, ok := queryParams["kpi_name"]; ok {
		query = query.Where("kpi_name = ?", kpiName)
	}

	if assignType, ok := queryParams["assign_type"]; ok {
		query = query.Where("assign_type = ?", assignType)
	}

	if rolesApplicable, ok := queryParams["roles_applicable"]; ok {
		query = query.Where("roles_applicable LIKE ?", "%"+rolesApplicable+"%")
	}

	if err := query.Find(&kpis).Error; err != nil {
		return kpis, err
	}

	return kpis, nil
}

func (r *KPIController) GetAllKPI(c *gin.Context) {
	queryParams := make(map[string]string)

	if kpiName := c.Query("kpi_name"); kpiName != "" {
		queryParams["kpi_name"] = kpiName
	}

	if assignType := c.Query("assign_type"); assignType != "" {
		queryParams["assign_type"] = assignType
	}

	if rolesApplicable := c.Query("roles_applicable"); rolesApplicable != "" {
		queryParams["roles_applicable"] = rolesApplicable
	}

	kpis, err := GetAllKPI(r.Db, queryParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, kpis)
}

// GetKPI retrieves a single KPI by ID
func GetKPIByID(db *gorm.DB, id uint) (models.KPI, error) {
	var kpi models.KPI

	if err := db.Preload("Measured").Preload("Questionaire").First(&kpi, id).Error; err != nil {
		return kpi, err
	}

	return kpi, nil
}

func (r *KPIController) GetKPIByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	kpi, err := GetKPIByID(r.Db, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, kpi)
}

func UpdateKPI(db *gorm.DB, id uint, updatedKPI models.KPI) (models.KPI, error) {
	var kpi models.KPI

	if err := db.Model(&kpi).Where("id = ?", id).Updates(updatedKPI).Error; err != nil {
		return kpi, err
	}

	if err := db.Model(&models.MeasuredData{}).Where("kpi_id = ?", id).Updates(models.MeasuredData{
		Key:   updatedKPI.Measured.Key,
		Value: updatedKPI.Measured.Value,
	}).Error; err != nil {
		return kpi, err
	}

	if err := db.Model(&models.QuestionaireData{}).Where("kpi_id = ?", id).Updates(models.QuestionaireData{
		Key:   updatedKPI.Questionaire.Key,
		Value: updatedKPI.Questionaire.Value,
	}).Error; err != nil {
		return kpi, err
	}

	if err := db.Preload("Measured").Preload("Questionaire").First(&kpi, id).Error; err != nil {
		return kpi, err
	}

	return kpi, nil
}

func (r *KPIController) UpdateKPI(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var updatedKPI models.KPI
	err = c.ShouldBindJSON(&updatedKPI)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	kpi, err := UpdateKPI(r.Db, uint(id), updatedKPI)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, kpi)
}

func DeleteKPI(db *gorm.DB, id uint) error {
    var kpi models.KPI
    if err := db.First(&kpi, id).Error; err != nil {
        return err
    }

    if err := db.Where("kpi_id = ?", id).Delete(&models.MeasuredData{}, "kpi_id").Error; err != nil {
        return err
    }
    if err := db.Where("kpi_id = ?", id).Delete(&models.QuestionaireData{}, "kpi_id").Error; err != nil {
        return err
    }

    if err := db.Delete(&kpi).Error; err != nil {
        return err
    }

    return nil
}




func (r *KPIController) DeleteKPI(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := DeleteKPI(r.Db, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
