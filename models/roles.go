package models

type RoleName string

const (
	Management RoleName = "Management"
	Supervisor RoleName = "Supervisor"
	HR         RoleName = "HR"
	Employee   RoleName = "Employee"
)

type Role struct {
	ID       uint     `gorm:"primaryKey" json:"role_id"`
	RoleName RoleName `gorm:"size:100; not null;" json:"role_name"`
	IsActive bool     `gorm:"not null" json:"is_active"`
}
