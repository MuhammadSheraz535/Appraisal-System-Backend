package service

import (
	"errors"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mrehanabbasi/appraisal-system-backend/controller"
	"github.com/mrehanabbasi/appraisal-system-backend/database"
	log "github.com/mrehanabbasi/appraisal-system-backend/logger"
	"github.com/mrehanabbasi/appraisal-system-backend/models"
	"gorm.io/gorm"
)

type EmployeeService struct {
	Db *gorm.DB
}

func NewEmployeeService() *EmployeeService {
	db := database.DB
	err := db.AutoMigrate(&models.Employee{}, &models.Role{})
	if err != nil {
		panic(err)
	}
	return &EmployeeService{Db: db}
}

// create employee

func (ec *EmployeeService) CreateEmployee(c *gin.Context) {
	log.Info("Initializing CreateEmployee handler function...")
	var employee models.Employee

	err := c.ShouldBindJSON(&employee)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if employee.Role == "" {
		employee.Role = "Employee"

	}
	// If a role is provided in the request, check if it exists in the DB and assign RoleId from the database
	if roleName := employee.Role; roleName != "" {
		roleId, err := controller.GetRoleIdFromDb(ec.Db, roleName)
		if err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		employee.RoleID = roleId
	}
	// checking supervisor exist in employee table
	if supID := employee.SupervisorID; supID != 0 {
		err := controller.ChecKSupervisorExist(ec.Db, supID)
		if err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
	}

	err = controller.CreateEmployee(ec.Db, &employee)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, employee)
}

// get all employee
func (uc *EmployeeService) GetEmployees(c *gin.Context) {
	log.Info("Initializing GetEmployees handler function...")
	var employees []models.Employee
	name := c.Query("name")
	role := c.Query("role")
	err := controller.GetEmployees(uc.Db, name, role, &employees)

	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, employees)
}

// get employee by id
func (ec *EmployeeService) GetEmployee(c *gin.Context) {
	log.Info("Initializing GetEmployee By ID handler function...")
	id, _ := strconv.Atoi(c.Param("id"))
	var employee models.Employee
	err := controller.GetEmployee(ec.Db, &employee, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("No employee found against the provided id")
			c.JSON(http.StatusNotFound, gin.H{"error": "No employee found against the provided id"})
			return
		}
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, employee)
}

// update Employee
func (ec *EmployeeService) UpdateEmployee(c *gin.Context) {
	log.Info("Initializing UpdateEmployee handler function...")
	var employee models.Employee
	id, _ := strconv.Atoi(c.Param("id"))
	err := controller.GetEmployee(ec.Db, &employee, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("No employee found for provided id")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No employee found for provided id"})
			return
		}

	}
	err = c.ShouldBindJSON(&employee)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// If a role is provided in the request, check if it exists in the DB and assign RoleId from the database
	if roleName := employee.Role; roleName != "" {
		roleId, err := controller.GetRoleIdFromDb(ec.Db, roleName)
		if err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		employee.RoleID = roleId
	}
	err = controller.UpdateEmployee(ec.Db, &employee)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, employee)
}

// delete Employee
func (ec *EmployeeService) DeleteEmployee(c *gin.Context) {
	log.Info("Initializing DeleteEmployee handler function...")
	var employee models.Employee
	id, _ := strconv.Atoi(c.Param("id"))
	_, err := controller.DeleteEmployee(ec.Db, &employee, id)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}
