package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

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

func GetAllRoles(db *gorm.DB, Role *[]models.Role) (err error) {
	err = db.Find(&Role).Error
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
		err = r.Db.Where("role_name = ? AND is_active = ?", roleName, isActive).Find(&role).Error
	} else if roleName != "" {
		err = r.Db.Where("role_name = ?", roleName).Find(&role).Error
	} else if isActive != "" {
		err = r.Db.Where("is_active = ?", isActive).Find(&role).Error
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
	err = db.Where("id = ?", id).First(&role).Error
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

func CreateRole(db *gorm.DB, role models.Role) (err error) {
	err = db.Create(&role).Error
	if err != nil {
		return err
	}
	return nil
}

// func (r *RoleController) CreateRole(c *gin.Context) {
//     var role models.Role
//     err := c.ShouldBindJSON(&role)
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }
//     if role.RoleName != "Management" && role.RoleName != "Supervisor" && role.RoleName != "HR" && role.RoleName != "Employee" {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid RoleName"})
//         return
//     }
//     err = CreateRole(r.Db, role)
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }
//     c.JSON(http.StatusOK, role)
// }

func (r *RoleController) CreateRole(c *gin.Context) {
	var role models.Role
	err := c.ShouldBindJSON(&role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	roleName := strings.ToLower(string(role.RoleName))
	if roleName != "management" && roleName != "supervisor" && roleName != "hr" && roleName != "employee" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid RoleName"})
		return
	}
	err = CreateRole(r.Db, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, role)
}

func UpdateRole(db *gorm.DB, role *models.Role) (err error) {
	err = db.Save(&role).Error
	return err
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
	c.ShouldBindJSON(&role)
	err = UpdateRole(r.Db, &role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, role)
}

func DeleteRole(db *gorm.DB, role *models.Role, id int) (int64, error) {
	db.Where("id = ?", id).Delete(&role)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (r *RoleController) DeleteRole(c *gin.Context) {
	var role models.Role
	id, _ := strconv.Atoi(c.Param("id"))
	_, err := DeleteRole(r.Db, &role, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
