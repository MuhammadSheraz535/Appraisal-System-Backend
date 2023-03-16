package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type StringSlice []string

func (s StringSlice) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *StringSlice) Scan(src interface{}) error {
	if src == nil {
		*s = nil
		return nil
	}

	var data []byte
	switch src.(type){
	case []byte:
		data = src.([]byte)
	case string:
		data = []byte(src.(string))
	default:
		return fmt.Errorf("invalid type for StringSlice: %T", src)
	}

	err := json.Unmarshal(data, &s)
	if err != nil {
		return fmt.Errorf("error unmarshaling StringSlice: %w", err)
	}

	return nil
}

type KPI struct {
	ID              uint               `gorm:"primaryKey" json:"kpi_id"`
	KPIName         string             `gorm:"size:100;not null" json:"kpi_name"`
	AssignType      StringSlice        `gorm:"type:json;not null" json:"assign_type"`
	Measured        []MeasuredData     `gorm:"-" json:"measured_data,omitempty"`
	Observatory     string             `gorm:"not null" json:"obs_data,omitempty"`
	Questionaire    []QuestionaireData `gorm:"-" json:"questionaire_data,omitempty"`
	Feedback        string             `gorm:"not null" json:"feedback_data,omitempty"`
	RolesApplicable StringSlice        `gorm:"type:json;not null" json:"roles_applicable,omitempty"`
}

type MeasuredData struct {
	Key   string      `json:"key"`
	Value int `json:"value"`
}

type QuestionaireData struct {
	Key   string      `json:"key"`
	Value bool `json:"value"`
}
