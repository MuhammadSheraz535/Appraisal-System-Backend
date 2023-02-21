package models

type RoleName string

const (
	Management RoleName = "Management"
	Supervisor RoleName = "Supervisor"
	HR         RoleName = "HR"
	Employee   RoleName = "Employee"
)

type Role struct {
	ID       uint `gorm:"primaryKey"`
	RoleName RoleName
	IsActive bool
}
