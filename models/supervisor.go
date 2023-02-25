package models

type Employee struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
type Supervisor struct {
	Name  string `json:"name" binding:"required,min=3,max=60"`
	Email string `json:"email" binding:"required,email"`
}
