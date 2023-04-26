package models

import (
	"github.com/lib/pq"
)

type KpiType struct {
	CommonModel
	KpiType      string `gorm:"not null;unique" json:"kpi_type"`
	BasicKpiType string `gorm:"not null" json:"basic_kpi_type"`
}

type AssignType struct {
	CommonModel
	AssignTypeId uint64 `gorm:"not null;unique" json:"assign_type_id"`
	AssignType   string `gorm:"not null;unique" json:"assign_type"`
}

type Kpi struct {
	CommonModel
	KpiName       string         `gorm:"size:100;not null;unique" json:"kpi_name"`
	AssignTypeID  uint64         `gorm:"not null" json:"assign_type"`
	AssignType    AssignType     `gorm:"references:AssignTypeId" json:"-"`
	KpiTypeID     string         `gorm:"not null" json:"kpi_type"`
	KpiType       KpiType        `gorm:"references:KpiType" json:"-"`
	ApplicableFor pq.StringArray `gorm:"type:text[];not null" json:"applicable_for"`
	Statement     string         `json:"statement,omitempty"`
}

type MultiKpi struct {
	CommonModel
	KpiName       string                  `json:"kpi_name"`
	AssignType    uint64                  `json:"assign_type"`
	KpiType       string                  `json:"kpi_type"`
	ApplicableFor pq.StringArray          `json:"applicable_for"`
	Statements    []MultiStatementKpiData `json:"statements,omitempty"`
}

type MultiStatementKpiData struct {
	CommonModel
	KpiID         uint64 `gorm:"foreignKey:ID" json:"-"`
	Statement     string `gorm:"not null" json:"statement"`
	CorrectAnswer string `gorm:"not null" json:"correct_answer"`
	Weightage     uint64 `gorm:"not null" json:"weightage"`
}
