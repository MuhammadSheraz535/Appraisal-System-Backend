package models

type Employee struct {
	Common
	Name         string `json:"name" gorm:"size:255;not null" binding:"required"`
	Email        string `json:"email" gorm:"unique;not null" binding:"required"`
	Role         string `json:"role_name"`
	RoleID       uint   `json:"role_id" gorm:"foreignKey:RoleID"`
	SupervisorID uint   `json:"supervisor_id,omitempty" gorm:"foreignKey:EmployeeID"`
}