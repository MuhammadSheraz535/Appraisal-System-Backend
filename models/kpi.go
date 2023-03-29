package models

import (
	"database/sql/driver"
	"errors"
	"strings"
)

type Questionnaire []string

type KpiType struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	KpiType string `gorm:"not null,unique" json:"kpi_type"`
}

type AssignType struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	AssignTypeId uint   `gorm:"not null,unique" json:"assign_type_id"`
	AssignType   string `gorm:"not null,unique" json:"assign_type"`
}

type Kpis struct {
	ID            uint   `gorm:"primaryKey" json:"kpi_id"`
	KpiName       string `gorm:"size:100;not null;unique" json:"kpi_name"`
	AssignType    uint64 `gorm:"not null" json:"assign_type"`
	KpiType       string `gorm:"not null" json:"kpi_type"`
	ApplicableFor string `gorm:"not null" json:"applicable_for"`
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

type ReqQuestionnaireKpi struct {
	KpisID        uint          `json:"-"`
	ID            uint          `json:"questionnaire_id"`
	Questionnaire Questionnaire `json:"questionnaire_data"`
}

type QuestionnaireKpi struct {
	KpisID        uint   `gorm:"foreignKey:KpisID" json:"-"`
	ID            uint   `gorm:"primaryKey" json:"questionnaire_id"`
	Questionnaire string `gorm:"type:VARCHAR(255)" json:"questionnaire_data"`
}

type MeasuredKpi struct {
	KpisID   uint   `gorm:"foreignKey:KpisID" json:"-"`
	ID       uint   `gorm:"primaryKey" json:"measured_id"`
	Measured string `gorm:"not null" json:"measured_data,omitempty"`
}

type ReqFeedBack struct {
	ID            uint   `json:"kpi_id"`
	KpiName       string `json:"kpi_name"`
	AssignType    uint64 `json:"assign_type"`
	KpiType       string `json:"kpi_type"`
	ApplicableFor string `json:"applicable_for"`
	Feedback      string `json:"feedback_data"`
}

type ReqObservatory struct {
	ID            uint   `json:"kpi_id"`
	KpiName       string `json:"kpi_name"`
	AssignType    uint64 `json:"assign_type"`
	KpiType       string `json:"kpi_type"`
	ApplicableFor string `json:"applicable_for"`
	Observatory   string `json:"obs_data"`
}

type ReqMeasured struct {
	ID            uint   `json:"kpi_id"`
	KpiName       string `json:"kpi_name"`
	AssignType    uint64 `json:"assign_type"`
	KpiType       string `json:"kpi_type"`
	ApplicableFor string `json:"applicable_for"`
	Measured      string `json:"measured_data"`
}

type ReqQuestionnaire struct {
	ID            uint     `json:"kpi_id"`
	KpiName       string   `json:"kpi_name"`
	AssignType    uint64   `json:"assign_type"`
	KpiType       string   `json:"kpi_type"`
	ApplicableFor string   `json:"applicable_for"`
	Questionnaire []string `json:"questionnaire_data"`
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
