package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrehanabbasi/appraisal-system-backend/database"
	"github.com/mrehanabbasi/appraisal-system-backend/models"
	"gorm.io/gorm"
)

type SupervisorController struct {
	db *gorm.DB
}

func NewSupervisorController() *SupervisorController {
	db := database.DB
	db.AutoMigrate(&models.Employee{})
	return &SupervisorController{db: db}
}

// Get Supervisors from Employee Tab

func (sc *SupervisorController) ConvertSupervisorToEmployee(c *gin.Context) {
	// Get the supervisor data from the request body
	var req models.SupervisorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Check if the email already exists in the employees table
	var existingSupervisor models.Employee
	if err := sc.db.Where("email = ?", req.Email).First(&existingSupervisor).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
		return
	}
	// Create a new Employee object with role set as "supervisor"
	employee := models.Employee{
		Name:  req.Name,
		Email: req.Email,
		Role:  "supervisor",
	}

	// Save the employee to the employees table
	if err := sc.db.Create(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set the role to "supervisor" in the response
	employee.Role = "supervisor"

	// Return the created employee
	c.JSON(http.StatusCreated, employee)
}

// GET Supervisors from employee table(with query parameter)
func (sc *SupervisorController) GetSupervisors(c *gin.Context) {
	name := c.Query("name")

	var supervisors []models.Employee
	if name != "" {
		// Get supervisors that match the specified name
		if err := sc.db.Where("role = ? AND name LIKE ?", "supervisor", "%"+name+"%").Find(&supervisors).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		// Get all supervisors
		if err := sc.db.Where("role = ?", "supervisor").Find(&supervisors).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Return the list of supervisors
	c.JSON(http.StatusOK, supervisors)
}

// Get Supervisor by ID
func (sc *SupervisorController) GetSupervisorById(c *gin.Context) {
	// Get the supervisor ID from the request parameters
	supervisorId := c.Param("id")

	// Query the employees table for an employee with the specified ID and role "supervisor"
	var supervisor models.Employee
	if err := sc.db.Where("id = ? AND role = ?", supervisorId, "supervisor").First(&supervisor).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "supervisor not found"})
		return
	}

	// Return the supervisor
	c.JSON(http.StatusOK, supervisor)
}

// Update Supervisor by ID
func (sc *SupervisorController) UpdateSupervisor(c *gin.Context) {
	// Get the supervisor ID from the request parameters
	supervisorId := c.Param("id")

	// Query the employees table for an employee with the specified ID and role "supervisor"
	var supervisor models.Employee
	if err := sc.db.Where("id = ? AND role = ?", supervisorId, "supervisor").First(&supervisor).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "supervisor not found"})
		return
	}

	// Get the updated supervisor data from the request body
	var req models.SupervisorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the supervisor's name and email in the employees table
	supervisor.Name = req.Name
	supervisor.Email = req.Email
	if err := sc.db.Save(&supervisor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the updated supervisor
	c.JSON(http.StatusOK, supervisor)
}
func (sc *SupervisorController) DeleteSupervisor(c *gin.Context) {
	// Get the supervisor ID from the request parameters
	supervisorId := c.Param("id")

	// Delete the supervisor with the specified ID from the employees table
	if err := sc.db.Where("id = ? AND role = ?", supervisorId, "supervisor").Delete(&models.Employee{}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "supervisor not found"})
		return
	}

	// Return a success response
	c.JSON(http.StatusOK, gin.H{"message": "supervisor deleted successfully"})
}
