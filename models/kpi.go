package models

import (
	"github.com/lib/pq"
)

type KpiType struct {
	ID           uint64 `gorm:"primaryKey" json:"id"`
	KpiType      string `gorm:"not null,unique" json:"kpi_type"`
	BasicKpiType string `gorm:"not null,unique" json:"basic_kpi_type"`
}

type AssignType struct {
	ID           uint64 `gorm:"primaryKey" json:"id"`
	AssignTypeId uint64 `gorm:"not null,unique" json:"assign_type_id"`
	AssignType   string `gorm:"not null,unique" json:"assign_type"`
}

type Kpi struct {
	ID            uint           `gorm:"primaryKey" json:"kpi_id"`
	KpiName       string         `gorm:"size:100;not null;unique" json:"kpi_name"`
	AssignType    uint64         `gorm:"not null" json:"assign_type"`
	KpiType       string         `gorm:"not null" json:"kpi_type"`
	ApplicableFor pq.StringArray `gorm:"type:text[];not null" json:"applicable_for"`
	Statement     string         `json:"statement,omitempty"`
}

type MultiKpi struct {
	ID            uint                    `json:"kpi_id"`
	KpiName       string                  `json:"kpi_name"`
	AssignType    uint64                  `json:"assign_type"`
	KpiType       string                  `json:"kpi_type"`
	ApplicableFor pq.StringArray          `json:"applicable_for"`
	Statements    []MultiStatementKpiData `json:"statements,omitempty"`
}

type MultiStatementKpiData struct {
	KpiID         uint   `gorm:"foreignKey:KpiID"`
	ID            uint64 `gorm:"primaryKey" json:"id"`
	Statement     string `gorm:"not null" json:"statement"`
	CorrectAnswer string `gorm:"not null" json:"correct_answer"`
	Weightage     uint64 `gorm:"not null" json:"weightage"`
}