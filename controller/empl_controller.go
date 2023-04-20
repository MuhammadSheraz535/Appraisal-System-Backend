package controller

import (
	"errors"

	"github.com/mrehanabbasi/appraisal-system-backend/models"
	"gorm.io/gorm"
)

// function to get the RoleId from the database based on the given role name
func GetRoleIdFromDb(db *gorm.DB, roleName string) (uint, error) {
	var role *models.Role
	if err := db.Table("roles").Where("role_name = ?", roleName).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, errors.New("invalid role specified")
		}
		return 0, err
	}

	return uint(role.ID), nil
}

func ChecKSupervisorExist(db *gorm.DB, id uint) error {
	var supervisor models.Employee
	if err := db.Table("employees").Where("id = ? AND role = ?", id, "supervisor").First(&supervisor).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("supervisor not found")
		}
		return err
	}
	return nil
}

// create a employee
func CreateEmployee(db *gorm.DB, employee *models.Employee) (err error) {

	if db.Table("employees").Where("email = ?", employee.Email).Find(&employee).RowsAffected > 0 {
		return errors.New("email is already registered")
	}

	if err = db.Table("employees").Create(&employee).Error; err != nil {
		return err
	}
	return nil
}

// get all employee
func GetEmployees(db *gorm.DB, name, role string, employees *[]models.Employee) (err error) {

	if name != "" && role != "" {
		err = db.Table("employees").Where("name = ? AND role = ?", name, role).Find(&employees).Error
		return err
	} else if name != "" {
		err = db.Table("employees").Where("name LIKE ?", "%"+name+"%").Find(&employees).Error
		return err

	} else if role != "" {
		err = db.Table("employees").Where("role LIKE ?", "%"+role+"%").Find(&employees).Error
		return err

	} else {
		err = db.Table("employees").Find(&employees).Error
		if err != nil {
			return err
		}
	}

	return nil

}

// get employee by id
func GetEmployee(db *gorm.DB, Employee *models.Employee, id int) (err error) {
	err = db.Table("employees").Where("id = ?", id).First(&Employee).Error
	if err != nil {
		return err
	}
	return nil
}

// update Employee
func UpdateEmployee(db *gorm.DB, Employee *models.Employee) (err error) {
	err = db.Model(&Employee).Updates(&Employee).Save(&Employee).Error
	return err
}

// delete Employee
func DeleteEmployee(db *gorm.DB, Employee *models.Employee, id int) (int64, error) {
	db = db.Table("employees").Debug().Model(&Employee).Where("id = ?", id).Take(&Employee).Unscoped().Delete(&Employee)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil

}
