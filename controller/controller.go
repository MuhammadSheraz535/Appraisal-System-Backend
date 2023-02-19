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

type SupervisorController struct {
	Db *gorm.DB
}

func New() *SupervisorController {
	db := database.Connect()
	db.AutoMigrate(&models.Employee{})
	return &SupervisorController{Db: db}
}

// create Supervisor Employee
func (sc *SupervisorController) CreateUser(c *gin.Context) {
	var user models.Employee
	c.BindJSON(&user)
	err := models.CreateUser(sc.Db, &user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, user)
}

// get All Suppervisors
func (sc *SupervisorController) GetUsers(c *gin.Context) {
	var user []models.Employee
	err := models.GetUsers(sc.Db, &user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, user)
}

// get Supervisor Employee by id
func (sc *SupervisorController) GetUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var user models.Employee
	err := models.GetUser(sc.Db, &user, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, user)
}

// update Supervisor Employee user
func (sc *SupervisorController) UpdateUser(c *gin.Context) {
	var user models.Employee
	id, _ := strconv.Atoi(c.Param("id"))
	err := models.GetUser(sc.Db, &user, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.BindJSON(&user)
	err = models.UpdateUser(sc.Db, &user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, user)
}

// delete Supervisor User
func (sc *SupervisorController) DeleteUser(c *gin.Context) {
	var user models.Employee
	id, _ := strconv.Atoi(c.Param("id"))
	err := models.DeleteUser(sc.Db, &user, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
