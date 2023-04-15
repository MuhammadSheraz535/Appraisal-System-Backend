package controller

import (
	"errors"

	"github.com/mrehanabbasi/appraisal-system-backend/models"
	"gorm.io/gorm"
)

// get all roles
func GetAllRoles(db *gorm.DB, role *[]models.Role) (err error) {
	err = db.Table("roles").Find(&role).Error
	if err != nil {
		return err
	}
	return nil
}

// get roles by id
func GetRoleByID(db *gorm.DB, role *models.Role, id int) (err error) {
	err = db.Table("roles").Where("id = ?", id).First(&role).Error
	if err != nil {
		return err
	}
	return nil
}

// create role
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

// updating role
func UpdateRole(db *gorm.DB, role *models.Role) error {
	var count int64
	if err := db.Table("roles").Where("role_name=?", role.RoleName).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("role name already exists")
	}

	var k models.Role
	if err := db.Table("roles").Updates(role).Error; err != nil {
		return err
	}
	role.CreatedAt = k.CreatedAt
	return nil
}

// delete roles
func DeleteRole(db *gorm.DB, role *models.Role, id int) error {
	err := db.Table("roles").Where("id = ?", id).Delete(&role).Error
	if err != nil {
		return err
	}
	return nil
}
