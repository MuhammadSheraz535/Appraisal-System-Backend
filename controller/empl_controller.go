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

func NewEmployeeController() *EmployeeController {
	db := database.DB
	db.AutoMigrate(&models.Employee{})
	return &EmployeeController{Db: db}
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

// create employee
func (ec *EmployeeController) CreateEmployee(c *gin.Context) {
	var employee models.Employee

	err := c.ShouldBindJSON(&employee)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if employee.Role == "" {
		employee.Role = "Employee"
	}

	err = CreateEmployee(ec.Db, &employee)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, employee)
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
		err = db.Table("employeea").Where("role LIKE ?", "%"+role+"%").Find(&employees).Error
		return err

	} else {
		err = db.Table("employees").Find(&employees).Error
		if err != nil {
			return err
		}
	}

	return nil

}

// get all employee
func (uc *EmployeeController) GetEmployees(c *gin.Context) {
	var employees []models.Employee
	name := c.Query("name")
	role := c.Query("role")
	err := GetEmployees(uc.Db, name, role, &employees)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, employees)
}

// get employee by id
func GetEmployee(db *gorm.DB, Employee *models.Employee, id int) (err error) {
	err = db.Table("employees").Where("id = ?", id).First(&Employee).Error
	if err != nil {
		return err
	}
	return nil
}

// get employee by id
func (ec *EmployeeController) GetEmployee(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var employee models.Employee
	err := GetEmployee(ec.Db, &employee, id)
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

// update Employee
func UpdateEmployee(db *gorm.DB, Employee *models.Employee) (err error) {
	err = db.Model(&Employee).Updates(&Employee).Save(&Employee).Error
	return err
}

// update Employee
func (ec *EmployeeController) UpdateEmployee(c *gin.Context) {
	var employee models.Employee
	id, _ := strconv.Atoi(c.Param("id"))
	err := GetEmployee(ec.Db, &employee, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Record not found"})
			return
		}

	}
	err = c.ShouldBindJSON(&employee)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = UpdateEmployee(ec.Db, &employee)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, employee)
}

// delete Employee
func DeleteEmployee(db *gorm.DB, Employee *models.Employee, id int) (int64, error) {
	db = db.Table("employees").Debug().Model(&Employee).Where("id = ?", id).Take(&Employee).Unscoped().Delete(&Employee)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil

}

// delete Employee
func (ec *EmployeeController) DeleteEmployee(c *gin.Context) {
	var employee models.Employee
	id, _ := strconv.Atoi(c.Param("id"))
	_, err := DeleteEmployee(ec.Db, &employee, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}
