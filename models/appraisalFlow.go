package models

type ApraisalFlow struct {
	Common
	FlowName  string     `json:"flow_name" binding:"required" gorm:"type:varchar(255);not null"`
	Createdby uint64     `json:"created_by" binding:"required" gorm:"not null"`
	IsActive  bool       `json:"is_active" binding:"required" gorm:"not null"`
	TeamId    uint64     `json:"team_id" binding:"required" gorm:"not null"`
	FlowStep  []FlowStep `json:"flowsteps" gorm:"foreignKey:FlowID;not null"`
}

type FlowStep struct {
	Common
	FlowID    uint64 `json:"-" binding:"required" gorm:"not null"`
	StepName  string `json:"step_name" binding:"required" gorm:"type:varchar(255);not null"`
	StepOrder uint64 `json:"step_order" binding:"required" gorm:"not null"`
	UserId    uint64 `json:"user_id" binding:"required" gorm:"not null"`
}
