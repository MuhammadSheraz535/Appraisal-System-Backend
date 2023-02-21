package models

type Employee struct {
	ID         uint       `json:"employee_id" gorm:"PrimaryKey"`
	Name       string     `gorm:"size:255;not null" json:"name"`
	Email      string     `gorm:"size:100;not null" json:"email"`
	Role       Role       `json:"role"  gorm:"foreignKey:ID"`
	Supervisor Supervisor `json:"supervisor" gorm:"foreignKey:ID;references:ID"`
}

type Role struct {
	ID       uint   `json:"role_id" gorm:"PrimaryKey"`
	Role     string `json:"employee_role"`
	IsActive bool   `json:"is_active"`
}
type Supervisor struct {
	ID         uint   `json:"supervisor_id" gorm:"PrimaryKey"`
	Name       string `json:"supervisor_name"`
	Email      string `json:"supervisor_email"`
	EmployeeID uint   `json:"-"`
}
