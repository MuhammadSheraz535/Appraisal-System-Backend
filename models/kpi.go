package models

import (
	"database/sql/driver"
	"errors"
	"strings"
)

type RolesApplicable []string

type KPI struct {
	ID              uint             `gorm:"primaryKey" json:"kpi_id"`
	KPIName         string           `gorm:"size:100;not null" json:"kpi_name"`
	AssignType      string           `gorm:"not null" json:"assign_type"`
	Observatory     string           `gorm:"not null" json:"obs_data,omitempty"`
	Feedback        string           `gorm:"not null" json:"feedback_data,omitempty"`
	Measured        MeasuredData     `gorm:"foreignKey:KPIID" json:"measured_data,omitempty"`
	Questionaire    QuestionaireData `gorm:"foreignKey:KPIID" json:"questionaire_data,omitempty"`
	RolesApplicable RolesApplicable  `gorm:"type:VARCHAR(255)" json:"roles_applicable"`
}

func (o *RolesApplicable) Scan(src any) error {
	bytes, ok := src.([]byte)
	if !ok {
		return errors.New("src value cannot cast to []byte")
	}
	*o = strings.Split(string(bytes), ",")
	return nil
}
func (o RolesApplicable) Value() (driver.Value, error) {
	if len(o) == 0 {
		return nil, nil
	}
	return strings.Join(o, ","), nil
}

type MeasuredData struct {
	KPIID uint `json:"kpi_id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type QuestionaireData struct {
	KPIID uint `json:"kpi_id"`
	Key   string `json:"key"`
	Value bool   `json:"value"`
}
