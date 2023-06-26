package service

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/mrehanabbasi/appraisal-system-backend/controller"
	"github.com/mrehanabbasi/appraisal-system-backend/database"
	log "github.com/mrehanabbasi/appraisal-system-backend/logger"
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
	log.Info("Initializing GetAllRoles handler function...")
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
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, role)
}

func (r *RoleService) GetRoleByID(c *gin.Context) {
	log.Info("Initializing GetRolesByID handler function...")
	id, _ := strconv.Atoi(c.Param("id"))
	var role models.Role
	err := controller.GetRoleByID(r.Db, &role, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("No Role found against the provided id")
			c.JSON(http.StatusNotFound, gin.H{"error": "No Role found against the provided id"})

		} else {
			log.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		return
	}
	c.JSON(http.StatusOK, role)
}

func (r *RoleService) CreateRole(c *gin.Context) {
	log.Info("Initializing CreateRole handler function...")
	var role models.Role
	err := c.ShouldBindJSON(&role)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err = controller.CreateRole(r.Db, role)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, role)
}

func (r *RoleService) UpdateRole(c *gin.Context) {
	log.Info("Initializing UpdateRoles handler function...")
	var role models.Role
	id, _ := strconv.Atoi(c.Param("id"))
	err := controller.GetRoleByID(r.Db, &role, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	err = c.ShouldBindJSON(&role)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = controller.UpdateRole(r.Db, &role)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, role)
}

func (r *RoleService) DeleteRole(c *gin.Context) {
	log.Info("Initializing DeleteRoles handler function...")
	var role models.Role
	id, _ := strconv.ParseUint(c.Param("id"), 10, 16)
	role.ID = uint16(id)
	err := controller.DeleteRole(r.Db, &role, role.ID)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
