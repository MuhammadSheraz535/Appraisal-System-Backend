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

// create user
func (uc *UserController) CreateUser(c *gin.Context) {
	var user models.Employee
	c.BindJSON(&user)
	err := models.CreateUser(uc.Db, &user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, user)
}

// get users
func (uc *UserController) GetUsers(c *gin.Context) {
	var user []models.Employee
	err := models.GetUsers(uc.Db, &user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, user)
}

// get user by id
func (uc *UserController) GetUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var user models.Employee
	err := models.GetUser(uc.Db, &user, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, user)
}


// get user by employee name
func (uc *UserController) GetUserByName(c *gin.Context) {
	name := c.Param("name")
	var user []models.Employee
	err := models.GetUserByName(uc.Db, &user, name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, user)
}
// get user by supervisor name
func (uc *UserController) GetSupervisorByName(c *gin.Context) {
	name := c.Param("name")
	var user []models.Supervisor
	err := models.GetSupervisorByName(uc.Db, &user, name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, user)
}
// get user by role
func (uc *UserController) GetByRole(c *gin.Context) {
	name := c.Param("name")
	var user []models.Role
	err := models.GetByRole(uc.Db, &user, name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, user)
}
// update user
func (uc *UserController) UpdateUser(c *gin.Context) {
	var user models.Employee
	id, _ := strconv.Atoi(c.Param("id"))
	err := models.GetUser(uc.Db, &user, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Record not found"})
			return
		}

	}
	c.BindJSON(&user)
	err = models.UpdateUser(uc.Db, &user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Record not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// delete user
func (uc *UserController) DeleteUser(c *gin.Context) {
	var user models.Employee
	id, _ := strconv.Atoi(c.Param("id"))
	_, err := models.DeleteUser(uc.Db, &user, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Record not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
