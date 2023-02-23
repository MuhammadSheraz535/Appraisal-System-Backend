package models

type Employee struct {
	ID    uint   `json:"employee_id" gorm:"PrimaryKey"`
	Name  string ` json:"name" gorm:"size:255;not null" binding:"required"`
	Email string `json:"email" gorm:"unique;not null" binding:"required"`
	Role  Role   `json:"role"  gorm:"foreignKey:ID"`
}

type Role struct {
	ID       uint   `json:"role_id" gorm:"PrimaryKey"`
	Role     string `json:"role"  binding:"required"`
	IsActive bool   `json:"is_active"  binding:"required"`
}
