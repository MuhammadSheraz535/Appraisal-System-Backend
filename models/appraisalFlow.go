package models

type AppraisalFlow struct {
	CommonModel
	FlowName  string     `json:"flow_name" gorm:"type:varchar(255);not null;unique"`
	CreatedBy uint64     `json:"created_by" gorm:"not null"`
	IsActive  bool       `json:"is_active" gorm:"not null"`
	TeamId    uint64     `json:"team_id" gorm:"not null"`
	FlowSteps []FlowStep `json:"flow_steps" gorm:"foreignKey:FlowID;not null"`
}

type FlowStep struct {
	CommonModel
	FlowID    uint64 `json:"-" gorm:"not null"`
	StepName  string `json:"step_name"  gorm:"type:varchar(255);not null"`
	StepOrder uint64 `json:"step_order"  gorm:"not null"`
	UserId    uint64 `json:"user_id"  gorm:"not null"`
}
