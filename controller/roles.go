package controller

import (
	"errors"
	"net/http"
	"strconv"

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

func GetAllRoles(db *gorm.DB, role *[]models.Role) (err error) {
	err = db.Table("roles").Find(&role).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *RoleController) GetAllRoles(c *gin.Context) {
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
		err = GetAllRoles(r.Db, &role)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, role)
}

func GetRoleByID(db *gorm.DB, role *models.Role, id int) (err error) {
	err = db.Table("roles").Where("id = ?", id).First(&role).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *RoleController) GetRoleByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var role models.Role
	err := GetRoleByID(r.Db, &role, id)
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

func CreateRole(db *gorm.DB, role models.Role) (models.Role, error) {
	var count int64
	if err := db.Table("roles").Where("role_name = ?", role.RoleName).Count(&count).Error; err != nil {
		return role, err
	}
	if count > 0 {
		return role, errors.New("role name already exists")
	}
	if err := db.Table("roles").Create(&role).Error; err != nil {
		return role, err
	}
	return role, nil
}

func (r *RoleController) CreateRole(c *gin.Context) {
	var role models.Role
	err := c.ShouldBindJSON(&role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	role, err = CreateRole(r.Db, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, role)
}

func UpdateRole(db *gorm.DB, role *models.Role) error {
	var count int64
	if err := db.Table("roles").Where("role_name=?", role.RoleName).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("role name already exists")
	}
	if err := db.Table("roles").Updates(role).Error; err != nil {
		return err
	}
	return nil
}

func (r *RoleController) UpdateRole(c *gin.Context) {
	var role models.Role
	id, _ := strconv.Atoi(c.Param("id"))
	err := GetRoleByID(r.Db, &role, id)
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

	err = UpdateRole(r.Db, &role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, role)
}

func DeleteRole(db *gorm.DB, role *models.Role, id int) error {
	err := db.Table("roles").Where("id = ?", id).Delete(&role).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *RoleController) DeleteRole(c *gin.Context) {
	var role models.Role
	id, _ := strconv.Atoi(c.Param("id"))
	role.ID = uint(id)
	err := DeleteRole(r.Db, &role, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
