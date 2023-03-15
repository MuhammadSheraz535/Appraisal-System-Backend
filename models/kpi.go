package models

type KPI struct {
	ID              uint               `gorm:"primaryKey" json:"kpi_id"`
	KPIName         string             `gorm:"size:100;not null" json:"kpi_name"`
	AssignType      []string           `gorm:"not null" json:"assign_type"`
	Measured        []MeasuredData     `gorm:"-" json:"measured_data,omitempty"`
	Observatory     string             `gorm:"not null" json:"obs_data,omitempty"`
	Questionaire    []QuestionaireData `gorm:"-" json:"questionaire_data,omitempty"`
	Feedback        string             `gorm:"not null" json:"feedback_data,omitempty"`
	RolesApplicable []string           `gorm:"not null" json:"roles_applicable,omitempty"`
}

type MeasuredData struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type QuestionaireData struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}
