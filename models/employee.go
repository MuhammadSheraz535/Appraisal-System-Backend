package models

type Employee struct {
	ID     uint   `gorm:"primary_key;auto_increment" json:"id"`
	Name   string `json:"name" gorm:"size:255;not null" binding:"required"`
	Email  string `json:"email" gorm:"unique;not null" binding:"required"`
	Role   string `json:"role"`
	RoleID uint   `json:"-" gorm:"foreignKey:RoleID"`
}
