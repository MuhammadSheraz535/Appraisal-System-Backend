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
	KpiName        string                  `gorm:"size:100;not null" json:"kpi_name"`
	KpiDescription string                  `gorm:"not null" json:"kpi_description"`
	AssignTypeID   uint64                  `gorm:"not null" json:"assign_type"`
	AssignType     AssignType              `gorm:"references:AssignTypeId;foreignKey:AssignTypeID" json:"-"`
	KpiTypeStr     string                  `gorm:"not null" json:"kpi_type"`
	KpiType        KpiType                 `gorm:"references:KpiType;foreignKey:KpiTypeStr" json:"-"`
	ApplicableFor  pq.StringArray          `gorm:"type:text[];not null" json:"applicable_for"`
	Statement      string                  `json:"statement,omitempty"`
	Statements     []MultiStatementKpiData `gorm:"foreignKey:KpiID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"statements,omitempty"`
}

type MultiStatementKpiData struct {
	CommonModel
	KpiID         uint64 `gorm:"not null" json:"-"`
	Statement     string `gorm:"not null" json:"statement"`
	CorrectAnswer string `gorm:"not null" json:"correct_answer"`
	Weightage     uint64 `gorm:"not null" json:"weightage"`
}
