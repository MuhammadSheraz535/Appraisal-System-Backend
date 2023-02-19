package models

import "gorm.io/gorm"

type Employee struct {
	gorm.Model
	E_ID  uint   `json:"employee_id" gorm:"PrimaryKey"`
	Name  string `gorm:"size:60;not null" json:"name"`
	Email string `gorm:"size:40;not null" json:"email"`
	S_ID  uint   `json:"supervisor_id" gorm:"foreignKey:E_ID;references:E_ID"`
}

//create Employee with Supervisor ID
func CreateUser(db *gorm.DB, User *Employee) (err error) {
	err = db.Create(User).Error
	if err != nil {
		return err
	}
	return nil
}

//get All Supervisor //Now it gives all employees
func GetUsers(db *gorm.DB, User *[]Employee) (err error) {
	err = db.Find(User).Error
	if err != nil {
		return err
	}
	return nil
}

//get Supervisor by id
func GetUser(db *gorm.DB, User *Employee, id int) (err error) {
	err = db.Where("id = ?", id).First(User).Error
	if err != nil {
		return err
	}
	return nil
}

//update Supervisor
func UpdateUser(db *gorm.DB, User *Employee) (err error) {
	db.Save(User)
	return nil
}

//delete Supervisor
func DeleteUser(db *gorm.DB, User *Employee, id int) (err error) {
	db.Where("id = ?", id).Delete(User)
	return nil
}
