package service

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/mrehanabbasi/appraisal-system-backend/controller"
	"github.com/mrehanabbasi/appraisal-system-backend/database"
	"github.com/mrehanabbasi/appraisal-system-backend/models"
)

type RoleService struct {
	Db *gorm.DB
}

func NewRoleService() *RoleService {
	db := database.DB
	err := db.AutoMigrate(&models.Role{})
	if err != nil {
		panic(err)
	}
	return &RoleService{Db: db}
}

func (r *RoleService) GetAllRoles(c *gin.Context) {
	var role []models.Role
	var err error

	roleName := c.Query("role_name")
	isActive := c.Query("is_active")

	if roleName != "" && isActive != "" {
		err = r.Db.Table("roles").Where("role_name = ? AND is_active = ?", roleName, isActive).Find(&role).Error
	} else if roleName != "" {
		err = r.Db.Table("roles").Where("role_name = ?", roleName).Find(&role).Error
	} else if isActive != "" {
		err = r.Db.Table("roles").Where("is_active = ?", isActive).Find(&role).Error
	} else {
		err = controller.GetAllRoles(r.Db, &role)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, role)
}

func (r *RoleService) GetRoleByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var role models.Role
	err := controller.GetRoleByID(r.Db, &role, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, role)
}

func (r *RoleService) CreateRole(c *gin.Context) {
	var role models.Role
	err := c.ShouldBindJSON(&role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err = controller.CreateRole(r.Db, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, role)
}

func (r *RoleService) UpdateRole(c *gin.Context) {
	var role models.Role
	id, _ := strconv.Atoi(c.Param("id"))
	err := controller.GetRoleByID(r.Db, &role, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	err = c.ShouldBindJSON(&role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = controller.UpdateRole(r.Db, &role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, role)
}

func (r *RoleService) DeleteRole(c *gin.Context) {
	var role models.Role
	id, _ := strconv.Atoi(c.Param("id"))
	role.ID = uint64(id)
	err := controller.DeleteRole(r.Db, &role, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
