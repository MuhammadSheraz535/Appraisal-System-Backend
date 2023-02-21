package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/mrehanabbasi/appraisal-system-backend/models"
)

type RoleController struct {
	db *gorm.DB
}

func NewRoleController(db *gorm.DB) *RoleController {
	return &RoleController{db}
}

func (r *RoleController) GetAllRoles(c *gin.Context) {
	var roles []models.Role

	// Apply filters if any
	roleName := c.Query("role_name")
	isActive := c.Query("is_active")
	query := r.db.Model(&models.Role{})
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

	if err := r.db.First(&role, c.Param("id")).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, role)
}

func (r *RoleController) CreateRole(c *gin.Context) {
	var role models.Role

	if err := c.ShouldBindJSON(&role); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := r.db.Create(&role).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, role)
}

func (r *RoleController) UpdateRole(c *gin.Context) {
	var role models.Role

	if err := r.db.First(&role, c.Param("id")).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err := c.ShouldBindJSON(&role); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := r.db.Save(&role).Error; err != nil {

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, role)
}

func (r *RoleController) DeleteRole(c *gin.Context) {
	if err := r.db.Delete(&models.Role{}, c.Param("id")).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusNoContent)
}
