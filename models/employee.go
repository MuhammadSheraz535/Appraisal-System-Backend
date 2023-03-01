package models

import "gorm.io/gorm"

type Employee struct {
	gorm.Model
	Name   string `json:"name" gorm:"size:255;not null" binding:"required"`
	Email  string `json:"email" gorm:"unique;not null" binding:"required"`
	Role   string `json:"role"`
	RoleID uint   `json:"-"`
}
