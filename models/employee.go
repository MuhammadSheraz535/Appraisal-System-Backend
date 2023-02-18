package models

import (
	"gorm.io/gorm"
)

type Role struct {
	ID       int    `json:"role_id" gorm:"PrimaryKey"`
	Role     string `json:"employee_role"`
	IsActive bool   `json:"is_active"`
}
type Employee struct {
	ID    uint   `json:"employee_id" gorm:"PrimaryKey"`
	Name  string `gorm:"size:255;not null" json:"name"`
	Email string `gorm:"size:100;not null" json:"email"`
	Role  Role   `json:"role"  gorm:"foreignKey:ID"`
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
	err = db.Preload("Role").Find(&User).Error
	if err != nil {
		return err
	}

	return nil

}

// get user by id
func GetUser(db *gorm.DB, User *Employee, id int) (err error) {
	err = db.Preload("Role").Where("id = ?", id).First(&User).Error
	if err != nil {
		return err
	}
	return nil
}

// update user
func UpdateUser(db *gorm.DB, User *Employee) (err error) {
	db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&User).Save(&User)
	return nil
}

// delete user
func DeleteUser(db *gorm.DB, User *Employee, id int) (int64, error) {
	db = db.Debug().Model(&User).Where("id = ?", id).Take(&User).Delete(&User)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
