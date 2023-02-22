package controller

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mrehanabbasi/appraisal-system-backend/database"
	"github.com/mrehanabbasi/appraisal-system-backend/models"
	"gorm.io/gorm"
)

type SupervisorController struct {
	Db *gorm.DB
}

func NewSupervisorController() *SupervisorController {
	db := database.DB
	db.AutoMigrate(&models.Supervisor{})
	return &SupervisorController{Db: db}
}

// create Supervisor type Employee
func CreateSupervisor(db *gorm.DB, supervisor *models.Supervisor) (err error) {
	var count int64
	db.Model(&models.Supervisor{}).Where("email=? OR s_id=?", supervisor.Email, supervisor.S_ID).Count(&count)
	if count > 0 {
		return errors.New("supervisor with the same email or supervisor ID already exists")
	}
	err = db.Create(supervisor).Error
	if err != nil {
		return err
	}
	return nil
}

// create Supervisor Employee
func (sc *SupervisorController) CreateSupervisor(c *gin.Context) {
	var supervisor models.Supervisor
	err := c.ShouldBindJSON(&supervisor)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	err = CreateSupervisor(sc.Db, &supervisor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, supervisor)
}

// get All Supervisor //Now it gives all employees
func GetSupervisors(db *gorm.DB, supervisor *[]models.Supervisor) (err error) {
	//Check before querying the database supervisor slice cannot be nil
	if supervisor == nil {
		return errors.New("supervisor slice pointer cannot be nil")
	}
	err = db.Find(supervisor).Error
	if err != nil {
		return err
	}
	return nil
}
func GetSupervisorsByName(db *gorm.DB, supervisor *[]models.Supervisor, name string) error {
	//check the white space provide instad of name of supervisor
	if strings.TrimSpace(name) == "" {
		return errors.New("name cannot be empty")
	}
	err := db.Where("name LIKE ?", "%"+name+"%").Find(supervisor).Error
	if err != nil {
		return err
	}
	return nil
}

func (sc *SupervisorController) GetSupervisors(c *gin.Context) {
	name := c.Query("_name")
	var supervisor []models.Supervisor

	if name != "" {
		err := GetSupervisorsByName(sc.Db, &supervisor, name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		err := GetSupervisors(sc.Db, &supervisor)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, supervisor)
}

// get supervisor by id
func GetSupervisorByID(db *gorm.DB, supervisor *models.Supervisor, id int) (err error) {
	err = db.Where("e_id = ?", id).First(&supervisor).Error
	if err != nil {
		return err
	}
	return nil
}

// get supervisor by id
func (uc *SupervisorController) GetSupervisorByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("e_id"))
	var supervisor models.Supervisor
	err := GetSupervisorByID(uc.Db, &supervisor, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, supervisor)
}

// update Supervisor
func UpdateSupervisor(db *gorm.DB, User *models.Supervisor) (err error) {
	err = db.Save(User).Error
	if err != nil {
		return err
	}

	return nil
}

// update Supervisor Employee user
func (sc *SupervisorController) UpdateSupervisor(c *gin.Context) {
	var user models.Supervisor
	id, _ := strconv.Atoi(c.Param("id"))
	err := GetSupervisorByID(sc.Db, &user, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.BindJSON(&user)
	err = UpdateSupervisor(sc.Db, &user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, user)
}

// delete Supervisor
func DeleteEmployee(db *gorm.DB, Employee *models.Supervisor, id int) (err error) {
	db = db.Debug().Model(&Employee).Where("id = ?", id).Take(&Employee).Delete(&Employee)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

// delete Supervisor
func (uc *SupervisorController) DeleteEmployee(c *gin.Context) {
	var supervisor models.Supervisor
	id, _ := strconv.Atoi(c.Param("id"))
	err := DeleteEmployee(uc.Db, &supervisor, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Record not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Supervisor deleted successfully"})
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
func (uc *SupervisorController) GetSupervisorByName(c *gin.Context) {
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
