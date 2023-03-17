package controller

import (
	"errors"

	"github.com/mrehanabbasi/appraisal-system-backend/models"
	"gorm.io/gorm"
)

const supervisorRoleName = "supervisor"

// Create Supervisor (Database)
func CreateSupervisor(db *gorm.DB, name string, email string, roleName string, roleID uint) (*models.Employee, error) {
	// Check if the email already exists in the employees table
	var existingEmployee models.Employee
	if err := db.Table("employees").Where("email = ?", email).First(&existingEmployee).Error; err == nil {
		return nil, errors.New("email already exists")
	}

	// Create a new Employee object
	employee := models.Employee{
		Name:   name,
		Email:  email,
		Role:   roleName,
		RoleID: roleID,
	}

	// Save the employee to the employees table
	if err := db.Table("employees").Create(&employee).Error; err != nil {
		return nil, err
	}

	return &employee, nil
}

// GET Supervisors from employee table(with query parameter) Database
func GetSupervisorsWithQuery(db *gorm.DB, name string) ([]models.Employee, error) {
	var supervisors []models.Employee
	if name != "" {
		// Get supervisors that match the specified name
		if err := db.Table("employees").Where("role = ? AND name LIKE ?", supervisorRoleName, "%"+name+"%").Find(&supervisors).Error; err != nil {
			return nil, err
		}
	} else {
		// Get all supervisors
		if err := db.Table("employees").Where("role = ?", supervisorRoleName).Find(&supervisors).Error; err != nil {
			return nil, err
		}
	}
	return supervisors, nil
}

// Get Supervisor by ID (Database)
func GetSupervisorByIdDB(db *gorm.DB, id string) (*models.Employee, error) {
	// Query the employees table for an employee with the specified ID and role "supervisor"
	var supervisor models.Employee
	if err := db.Table("employees").Where("id = ? AND role = ?", id, supervisorRoleName).First(&supervisor).Error; err != nil {
		return nil, err
	}

	return &supervisor, nil
}

// Update Supervisor by ID (Database)
func UpdateSupervisorInDatabase(db *gorm.DB, supervisorId string, req models.Supervisor) error {
	// Query the employees table for an employee with the specified ID and role "supervisor"
	var supervisor models.Employee
	if err := db.Table("employees").Where("id = ? AND role = ?", supervisorId, supervisorRoleName).First(&supervisor).Error; err != nil {
		return err
	}

	// Update the supervisor's name and email in the employees table
	supervisor.Name = req.Name
	supervisor.Email = req.Email
	if err := db.Table("employees").Save(&supervisor).Error; err != nil {
		return err
	}

	return nil
}

// Delete the supervisor from employee table(Database)
func DeleteSupervisorFromDB(db *gorm.DB, supervisorId string) error {
	// Delete the supervisor with the specified ID from the employees table
	if err := db.Table("employees").Where("id = ? AND role = ?", supervisorId, supervisorRoleName).Delete(&models.Employee{}).Error; err != nil {
		return err
	}

	return nil
}
