package models

type Role struct {
	Common
	RoleName string `gorm:"size:100;not null;unique" json:"role_name"`
	IsActive bool   `gorm:"not null" json:"is_active"`
}