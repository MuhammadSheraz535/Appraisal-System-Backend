package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/mrehanabbasi/appraisal-system-backend/database"
	"github.com/mrehanabbasi/appraisal-system-backend/models"
)

type RoleController struct {
	Db *gorm.DB
}

func NewRoleController() *RoleController {
	db := database.DB
	db.AutoMigrate(&models.Role{})
	return &RoleController{Db: db}
}

func (r *RoleController) GetAllRoles(c *gin.Context) {
	var roles []models.Role

	// Apply filters if any
	roleName := c.Query("role_name")
	isActive := c.Query("is_active")
	query := r.Db.Model(&models.Role{})
	if roleName != "" {
		query = query.Where("role_name = ?", roleName)
	}
	if isActive != "" {
		query = query.Where("is_active = ?", isActive)
	}

	if err := query.Find(&roles).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, roles)
}

func (r *RoleController) GetRoleByID(c *gin.Context) {
	var role models.Role

	if err := r.Db.First(&role, c.Param("role_id")).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, role)
}


func CreateRole(db *gorm.DB, role models.Role) (err error) {
	err = db.Create(&role).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *RoleController) CreateRole(c *gin.Context) {
	var role models.Role
	err := c.ShouldBindJSON(&role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = CreateRole(r.Db, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, role)
}

func (r *RoleController) UpdateRole(c *gin.Context) {
	var role models.Role

	if err := r.Db.First(&role, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.Db.Save(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, role)
}

func (r *RoleController) DeleteRole(c *gin.Context) {
	if err := r.Db.Delete(&models.Role{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
