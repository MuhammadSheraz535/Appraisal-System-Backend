package models

type RoleName string

const (
	Management RoleName = "Management"
	Supervisor RoleName = "Supervisor"
	HR         RoleName = "HR"
	Employee   RoleName = "Employee"
)

type Role struct {
	ID       uint     `gorm:"column:id;primaryKey" json:"id"`
	RoleName RoleName `gorm:"column:role_name"`
	IsActive bool     `gorm:"column:is_active"`
}
