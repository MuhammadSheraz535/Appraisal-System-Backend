package models

type RoleName string

const (
	ManagementRole RoleName = "Management"
	SupervisorRole RoleName = "Supervisor"
	HRRole         RoleName = "HR"
	EmployeeRole   RoleName = "Employee"
)

type Role struct {
	ID       uint     `gorm:"primaryKey" json:"role_id"`
	RoleName RoleName `gorm:"size:100; not null;unique" json:"role_name"`
	IsActive bool     `gorm:"not null" json:"is_active"`
}
