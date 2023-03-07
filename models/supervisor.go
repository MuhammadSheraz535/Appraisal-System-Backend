package models

type Supervisor struct {
	Name  string `json:"name" binding:"required,min=3,max=60"`
	Email string `json:"email" binding:"required,email"`
}
