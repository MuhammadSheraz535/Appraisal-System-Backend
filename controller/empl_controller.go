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

type UserController struct {
	Db *gorm.DB
}

func New() *UserController {
	db := database.Connect()
	db.AutoMigrate(&models.Employee{})
	return &UserController{Db: db}
}

// create a user
func CreateUser(db *gorm.DB, User *models.Employee) (err error) {
	err = db.Create(&User).Error
	if err != nil {
		return err
	}
	return nil
}

// get all users
func GetUsers(db *gorm.DB, User *[]models.Employee) (err error) {
	err = db.Preload("Role").Preload("Supervisor").Find(&User).Error
	if err != nil {
		return err
	}

	return nil

}

// get user by id
func GetUser(db *gorm.DB, User *models.Employee, id int) (err error) {
	err = db.Preload("Role").Preload("Supervisor").Where("id = ?", id).First(&User).Error
	if err != nil {
		return err
	}
	return nil
}

// get user by employee name
func GetUserByName(db *gorm.DB, User *[]models.Employee, name string) (err error) {

	err = db.Preload("Role").Preload("Supervisor").Where("name = ?", name).Statement.Find(&User).Error
	if err != nil {
		return err
	}
	return nil
}

// - _role_: Share employee list with the specified role.
func GetByRole(db *gorm.DB, User *[]models.Role, name string) (err error) {

	err = db.Where("role = ?", name).Find(&User).Error
	if err != nil {
		return err
	}
	return nil
}

// get user by supervisor name
func GetSupervisorByName(db *gorm.DB, User *[]models.Supervisor, name string) (err error) {

	err = db.Where("name = ?", name).Find(&User).Error
	if err != nil {
		return err
	}
	return nil
}

// update user
func UpdateUser(db *gorm.DB, User *models.Employee) (err error) {
	err = db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&User).Save(&User).Error
	return err
}

// delete user
func DeleteUser(db *gorm.DB, User *models.Employee, id int) (int64, error) {
	db = db.Debug().Model(&User).Where("id = ?", id).Take(&User).Delete(&User)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// create user
func (uc *UserController) CreateUser(c *gin.Context) {
	var user models.Employee
	c.ShouldBindJSON(&user)
	err := CreateUser(uc.Db, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// get users
func (uc *UserController) GetUsers(c *gin.Context) {
	var user []models.Employee
	err := GetUsers(uc.Db, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// get user by id
func (uc *UserController) GetUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var user models.Employee
	err := GetUser(uc.Db, &user, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// get user by employee name
func (uc *UserController) GetUserByName(c *gin.Context) {
	name := c.Param("name")
	var user []models.Employee
	err := GetUserByName(uc.Db, &user, name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// get user by supervisor name
func (uc *UserController) GetSupervisorByName(c *gin.Context) {
	name := c.Param("name")
	var user []models.Supervisor
	err := GetSupervisorByName(uc.Db, &user, name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// get user by role
func (uc *UserController) GetByRole(c *gin.Context) {
	name := c.Param("name")
	var user []models.Role
	err := GetByRole(uc.Db, &user, name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// update user
func (uc *UserController) UpdateUser(c *gin.Context) {
	var user models.Employee
	id, _ := strconv.Atoi(c.Param("id"))
	err := GetUser(uc.Db, &user, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Record not found"})
			return
		}

	}
	c.ShouldBindJSON(&user)
	err = UpdateUser(uc.Db, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Record not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// delete user
func (uc *UserController) DeleteUser(c *gin.Context) {
	var user models.Employee
	id, _ := strconv.Atoi(c.Param("id"))
	_, err := DeleteUser(uc.Db, &user, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Record not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
