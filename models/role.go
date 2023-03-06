package models

type RoleName string

type Role struct {
	ID       uint   `gorm:"primaryKey" json:"role_id"`
	RoleName string `gorm:"size:100; not null;unique" json:"role_name"`
	IsActive bool   `gorm:"not null" json:"is_active"`
}
