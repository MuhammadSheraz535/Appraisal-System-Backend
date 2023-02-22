package models

type Supervisor struct {
	ID    uint   `json:"employee_id" gorm:"PrimaryKey"`
	Name  string `gorm:"size:60;not null" json:"name"`
	Email string `gorm:"size:40;not null;unique" json:"email"`
	S_ID  uint   `json:"supervisor_id" gorm:"foreignKey:ID;references:ID"`
}
