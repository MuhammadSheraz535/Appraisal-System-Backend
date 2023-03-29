package models

import (
	"database/sql/driver"
	"errors"
	"strings"
)

type Questionnaire []string

type KpiType struct {
	ID         uint
	AssignType string `gorm:"not null,unique" json:"assign_type"`
}

type Kpis struct {
	ID         uint   `gorm:"primaryKey" json:"kpi_id"`
	KpiName    string `gorm:"size:100;not null;unique" json:"kpi_name"`
	AssignType string `gorm:"not null" json:"assign_type"`
	KpiType    string `gorm:"not null" json:"kpi_type"`
}

type FeedbackKpi struct {
	KpisID   uint   `gorm:"foreignKey:KpisID" json:"-"`
	ID       uint   `gorm:"primaryKey" json:"feedback_id"`
	FeedBack string `gorm:"not null" json:"feedback_data,omitempty"`
}

type ObservatoryKpi struct {
	KpisID      uint   `gorm:"foreignKey:KpisID" json:"-"`
	ID          uint   `gorm:"primaryKey" json:"observatory_id"`
	Observatory string `gorm:"not null" json:"obs_data,omitempty"`
}

type QuestionnaireKpi struct {
	KpisID        uint          `gorm:"foreignKey:KpisID" json:"-"`
	ID            uint          `gorm:"primaryKey" json:"questionnaire_id"`
	Questionnaire Questionnaire `gorm:"type:VARCHAR(255)" json:"questionnaire_data"`
}

type SingleQuestionnaireKpi struct {
	KpisID        uint   `gorm:"foreignKey:KpisID" json:"-"`
	ID            uint   `gorm:"primaryKey" json:"questionnaire_id"`
	Questionnaire string `gorm:"type:VARCHAR(255)" json:"questionnaire_data"`
}

type MeasuredKpi struct {
	KpisID   uint   `gorm:"foreignKey:KpisID" json:"-"`
	ID       uint   `gorm:"primaryKey" json:"measured_id"`
	Measured string `gorm:"not null" json:"measured_data,omitempty"`
}

func (o *Questionnaire) Scan(src any) error {
	bytes, ok := src.([]byte)
	if !ok {
		return errors.New("src value cannot cast to []byte")
	}
	*o = strings.Split(string(bytes), ",")
	return nil
}
func (o Questionnaire) Value() (driver.Value, error) {
	if len(o) == 0 {
		return nil, nil
	}
	return strings.Join(o, ","), nil
}
