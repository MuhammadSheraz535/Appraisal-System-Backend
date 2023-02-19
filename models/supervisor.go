package models

import "gorm.io/gorm"

type Employee struct {
	gorm.Model
	E_ID  uint   `json:"employee_id" gorm:"PrimaryKey"`
	Name  string `gorm:"size:60;not null" json:"name"`
	Email string `gorm:"size:40;not null" json:"email"`
	S_ID  uint   `json:"supervisor_id" gorm:"foreignKey:E_ID"`
}
