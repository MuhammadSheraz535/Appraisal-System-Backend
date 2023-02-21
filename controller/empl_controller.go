package controller

import (
	"errors"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mrehanabbasi/appraisal-system-backend/database"
	"github.com/mrehanabbasi/appraisal-system-backend/models"
	"gorm.io/gorm"
)

type EmployeeController struct {
	Db *gorm.DB
}

func New() *EmployeeController {
	db := database.DB
	db.AutoMigrate(&models.Employee{})
	return &EmployeeController{Db: db}
}

// create a employee
func CreateEmployee(db *gorm.DB, Employee *models.Employee) (err error) {
	err = db.Create(&Employee).Error
	if err != nil {
		return err
	}
	return nil
}

// create employee
func (uc *EmployeeController) CreateEmployee(c *gin.Context) {
	var employee models.Employee
	c.BindJSON(&employee)
	err := CreateEmployee(uc.Db, &employee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, employee)
}

// get all employee
func GetEmployees(db *gorm.DB, Employee *[]models.Employee) (err error) {
	err = db.Preload("Role").Preload("Supervisor").Find(&Employee).Error
	if err != nil {
		return err
	}

	return nil

}

// get all employee
func (uc *EmployeeController) GetEmployees(c *gin.Context) {
	var employee []models.Employee
	err := GetEmployees(uc.Db, &employee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, employee)
}

// get employee by id
func GetEmployee(db *gorm.DB, Employee *models.Employee, id int) (err error) {
	err = db.Preload("Role").Preload("Supervisor").Where("id = ?", id).First(&Employee).Error
	if err != nil {
		return err
	}
	return nil
}

// get employee by id
func (uc *EmployeeController) GetEmployee(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var employee models.Employee
	err := GetEmployee(uc.Db, &employee, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, employee)
}

// get employee list by employee name
func GetEmployeeByName(db *gorm.DB, Employee *[]models.Employee, name string) (err error) {

	err = db.Preload("Role").Preload("Supervisor").Where("name = ?", name).Statement.Find(&Employee).Error
	if err != nil {
		return err
	}
	return nil
}

// get employee list by name
func (uc *EmployeeController) GetEmployeeByName(c *gin.Context) {
	name := c.Param("name")
	var employee []models.Employee
	err := GetEmployeeByName(uc.Db, &employee, name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, employee)
}

// get employee list  by role

func GetByRole(db *gorm.DB, USER *[]models.Role, name string) (err error) {

	err = db.Where("role = ?", name).Find(&USER).Error
	if err != nil {
		return err
	}
	return nil
}

// get employee list by role
func (uc *EmployeeController) GetByRole(c *gin.Context) {
	name := c.Param("name")
	var employee []models.Role
	err := GetByRole(uc.Db, &employee, name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, employee)
}

// get supervisor list by supervisor name
func GetSupervisorByName(db *gorm.DB, USER *[]models.Supervisor, name string) (err error) {

	err = db.Where("name = ?", name).Find(&USER).Error
	if err != nil {
		return err
	}
	return nil
}

// get supervisor list by supervisor name
func (uc *EmployeeController) GetSupervisorByName(c *gin.Context) {
	name := c.Param("name")
	var Supervisor []models.Supervisor
	err := GetSupervisorByName(uc.Db, &Supervisor, name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, Supervisor)
}

// update Employee
func UpdateEmployee(db *gorm.DB, Employee *models.Employee) (err error) {
	err = db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&Employee).Save(&Employee).Error
	return err
}

// update Employee
func (uc *EmployeeController) UpdateEmployee(c *gin.Context) {
	var employee models.Employee
	id, _ := strconv.Atoi(c.Param("id"))
	err := GetEmployee(uc.Db, &employee, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Record not found"})
			return
		}

	}
	c.ShouldBindJSON(&employee)
	err = UpdateEmployee(uc.Db, &employee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Record not found"})
		return
	}
	c.JSON(http.StatusOK, employee)
}

// delete Employee
func DeleteEmployee(db *gorm.DB, Employee *models.Employee, id int) (int64, error) {
	db = db.Debug().Model(&Employee).Where("id = ?", id).Take(&Employee).Delete(&Employee)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// delete Employee
func (uc *EmployeeController) DeleteEmployee(c *gin.Context) {
	var employee models.Employee
	id, _ := strconv.Atoi(c.Param("id"))
	_, err := DeleteEmployee(uc.Db, &employee, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Record not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
