package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Id       int `json:"id" gorm:"primary_key"`
	Name     int `json:"name"`
	Email    int `json:"email"`
	Password int `json:"password"`
}
