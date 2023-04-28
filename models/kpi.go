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
	KpiName        string                  `gorm:"size:100;not null" json:"kpi_name" binding:"required"`
	KpiDescription string                  `gorm:"not null" json:"kpi_description" binding:"required"`
	AssignTypeID   uint64                  `gorm:"not null" json:"assign_type" binding:"required"`
	AssignType     AssignType              `gorm:"references:AssignTypeId;foreignKey:AssignTypeID" json:"-"`
	KpiTypeStr     string                  `gorm:"not null" json:"kpi_type" binding:"required"`
	KpiType        KpiType                 `gorm:"references:KpiType;foreignKey:KpiTypeStr" json:"-"`
	ApplicableFor  pq.StringArray          `gorm:"type:text[];not null" json:"applicable_for" binding:"required"`
	Statement      string                  `json:"statement,omitempty"`
	Statements     []MultiStatementKpiData `gorm:"foreignKey:KpiID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"statements,omitempty"`
}

type MultiStatementKpiData struct {
	CommonModel
	KpiID         uint64 `gorm:"not null" json:"-"`
	Statement     string `gorm:"not null" json:"statement" binding:"required"`
	CorrectAnswer string `gorm:"not null" json:"correct_answer" binding:"required"`
	Weightage     uint64 `gorm:"not null" json:"weightage" binding:"required"`
}
