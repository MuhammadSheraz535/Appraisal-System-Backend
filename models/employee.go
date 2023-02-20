package models

import (
	"gorm.io/gorm"
)

type Employee struct {
	ID         uint       `json:"employee_id" gorm:"PrimaryKey"`
	Name       string     `gorm:"size:255;not null" json:"name"`
	Email      string     `gorm:"size:100;not null" json:"email"`
	Role       Role       `json:"role"  gorm:"foreignKey:ID"`
	Supervisor Supervisor `json:"supervisor" gorm:"foreignKey:ID;references:ID"`
}

type Role struct {
	ID       uint   `json:"role_id" gorm:"PrimaryKey"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}
type Supervisor struct {
	ID         uint   `json:"supervisor_id" gorm:"PrimaryKey"`
	Name       string `json:"name"`
	Email      string `json:"supervisor_email"`
	EmployeeID uint   `json:"-"`
}

// create a user
func CreateUser(db *gorm.DB, User *Employee) (err error) {
	err = db.Create(&User).Error
	if err != nil {
		return err
	}
	return nil
}

// get all users
func GetUsers(db *gorm.DB, User *[]Employee) (err error) {
	err = db.Preload("Role").Preload("Supervisor").Find(&User).Error
	if err != nil {
		return err
	}

	return nil

}

// get user by id
func GetUser(db *gorm.DB, User *Employee, id int) (err error) {
	err = db.Preload("Role").Preload("Supervisor").Where("id = ?", id).First(&User).Error
	if err != nil {
		return err
	}
	return nil
}

// get user by employee name
func GetUserByName(db *gorm.DB, User *[]Employee, name string) (err error) {

	err = db.Preload("Role").Preload("Supervisor").Where("name = ?", name).Statement.Find(&User).Error
	if err != nil {
		return err
	}
	return nil
}

// - _role_: Share employee list with the specified role.
func GetByRole(db *gorm.DB, User *[]Role, name string) (err error) {

	err = db.Where("role = ?", name).Find(&User).Error
	if err != nil {
		return err
	}
	return nil
}

// get user by supervisor name
func GetSupervisorByName(db *gorm.DB, User *[]Supervisor, name string) (err error) {

	err = db.Where("name = ?", name).Find(&User).Error
	if err != nil {
		return err
	}
	return nil
}

// update user
func UpdateUser(db *gorm.DB, User *Employee) (err error) {
	err = db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&User).Save(&User).Error
	return err
}

// delete user
func DeleteUser(db *gorm.DB, User *Employee, id int) (int64, error) {
	db = db.Debug().Model(&User).Where("id = ?", id).Take(&User).Delete(&User)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
