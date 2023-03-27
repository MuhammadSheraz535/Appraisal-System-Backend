package models

import (
	"database/sql/driver"
	"errors"
	"strings"
)


type Questionaire []string

type KpiType struct {
	ID         uint
	AssignType string `gorm:"not null,unique" json:"assign_type"`
}

type KpiCommon struct {
	ID         uint   `gorm:"primaryKey" json:"kpi_id"`
	KpiName    string `gorm:"size:100;not null,unique" json:"kpi_name"`
	AssignType string `gorm:"not null" json:"assign_type"`
	KpiType    string `gorm:"not null" json:"kpi_type"`
}

type FeedbackKpi struct {
	KpiCommon
	FeedBack    string    `gorm:"not null" json:"feedback_data,omitempty"`
}

type ObservatoryKpi struct {
	KpiCommon
	Observatory string    `gorm:"not null" json:"obs_data,omitempty"`
}

type QuestionaireKpi struct {
	KpiCommon
	Questionaire Questionaire `gorm:"type:VARCHAR(255)" json:"questionaire_data"`
}

type MeasuredKpi struct {
	KpiCommon
	Measured    string    `gorm:"not null" json:"measured_data,omitempty"`
}


func (o *Questionaire) Scan(src any) error {
	bytes, ok := src.([]byte)
	if !ok {
		return errors.New("src value cannot cast to []byte")
	}
	*o = strings.Split(string(bytes), ",")
	return nil
}
func (o Questionaire) Value() (driver.Value, error) {
	if len(o) == 0 {
		return nil, nil
	}
	return strings.Join(o, ","), nil
}
