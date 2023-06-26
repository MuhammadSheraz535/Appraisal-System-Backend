package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/mrehanabbasi/appraisal-system-backend/controller"
	"github.com/mrehanabbasi/appraisal-system-backend/database"
	log "github.com/mrehanabbasi/appraisal-system-backend/logger"
	"github.com/mrehanabbasi/appraisal-system-backend/models"
)

type SupervisorService struct {
	db *gorm.DB
}

func NewSupervisorService() *SupervisorService {
	db := database.DB
	err := db.AutoMigrate(&models.Employee{})
	if err != nil {
		panic(err)
	}
	return &SupervisorService{db: db}
}

const supervisorRoleName = "supervisor"

// Create Supervisor (handler)
func (sc *SupervisorService) ConvertSupervisorToEmployee(c *gin.Context) {
	log.Info("Initializing ConvertSupervisorToEmployee handler function...")
	// Get the supervisor data from the request body
	var req models.Supervisor
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the supervisor role from the roles table
	var supervisorRole models.Role
	if err := sc.db.Table("roles").Where("role_name = ?", supervisorRoleName).First(&supervisorRole).Error; err != nil {
		log.Error("supervisor role does not exist")
		c.JSON(http.StatusBadRequest, gin.H{"error": "supervisor role does not exist"})
		return
	}

	// Create a new employee with supervisor role
	employee, err := controller.CreateSupervisor(sc.db, req.Name, req.Email, supervisorRoleName, uint(supervisorRole.ID))
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the role to "supervisor" in the response
	employee.Role = supervisorRoleName

	// Return the created employee
	c.JSON(http.StatusCreated, employee)
}

// GET Supervisors from employee table(with query parameter) Handler
func (sc *SupervisorService) GetSupervisors(c *gin.Context) {
	name := c.Query("name")

	supervisors, err := controller.GetSupervisorsWithQuery(sc.db, name)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the list of supervisors
	c.JSON(http.StatusOK, supervisors)
}

// Get Supervisor by ID (handler)
func (sc *SupervisorService) GetSupervisorById(c *gin.Context) {
	log.Info("Initializing GetSupervisorById handler function...")
	// Get the supervisor ID from the request parameters
	supervisorId := c.Param("id")

	// Get the supervisor from the database
	supervisor, err := controller.GetSupervisorByIdDB(sc.db, supervisorId)
	if err != nil {
		log.Error("supervisor not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "supervisor not found"})
		return
	}

	// Return the supervisor
	c.JSON(http.StatusOK, *supervisor)
}

// Update Supervisor by ID (Handler)
func (sc *SupervisorService) UpdateSupervisor(c *gin.Context) {
	log.Info("Initializing UpdateSupervisor handler function...")
	// Get the supervisor ID from the request parameters
	supervisorId := c.Param("id")

	// Get the updated supervisor data from the request body
	var req models.Supervisor
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the supervisor in the database
	if err := controller.UpdateSupervisorInDatabase(sc.db, supervisorId, req); err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Error("supervisor not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "supervisor not found"})
			return
		}
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Query the employees table for the updated supervisor
	var updatedSupervisor models.Employee
	if err := sc.db.Table("employees").Where("id = ?", supervisorId).First(&updatedSupervisor).Error; err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the updated supervisor
	c.JSON(http.StatusOK, updatedSupervisor)
}

// DeleteSupervisor deletes a supervisor from the employee table (handler)
func (sc *SupervisorService) DeleteSupervisor(c *gin.Context) {
	log.Info("Initializing DeleteSupervisor handler function...")
	// Get the supervisor ID from the request parameters
	supervisorId := c.Param("id")

	// Call the database function to delete the supervisor
	if err := controller.DeleteSupervisorFromDB(sc.db, supervisorId); err != nil {
		log.Error("supervisor not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "supervisor not found"})
		return
	}

	// Return a success response
	c.Status(http.StatusNoContent)
}
