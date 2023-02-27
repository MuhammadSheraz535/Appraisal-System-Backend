package models

import "gorm.io/gorm"

type Employee struct {
	gorm.Model
	Name   string ` json:"name" gorm:"size:255;not null" binding:"required"`
	Email  string `json:"email" gorm:"unique;not null" binding:"required"`
	Role   Role   `json:"role"`
	RoleID uint `json:"-"`
}

type Role struct {
	ID       int   `json:"role_id" gorm:"PrimaryKey"`
	Role     string `json:"role"  binding:"required"`
	IsActive bool   `json:"is_active"`
	
}
