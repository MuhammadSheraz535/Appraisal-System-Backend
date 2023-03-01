package models

type Employee struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Name   string ` json:"name" gorm:"size:255;not null" binding:"required"`
	Email  string `json:"email" gorm:"unique;not null" binding:"required"`
	RoleID uint   `json:"roleid"`
	Role   string `json:"role" gorm:"foreignKey:RoleID"`
}
