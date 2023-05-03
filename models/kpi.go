package models

import (
	"github.com/lib/pq"
	"github.com/mrehanabbasi/appraisal-system-backend/constants"
)

type BasicKpiType string
type AssignTypeStr string

type KpiType struct {
	CommonModel
	KpiType      string       `gorm:"not null;unique" json:"kpi_type"`
	BasicKpiType BasicKpiType `gorm:"not null" json:"basic_kpi_type" binding:"enum"`
}

type AssignType struct {
	CommonModel
	AssignTypeId uint64        `gorm:"not null;unique" json:"assign_type_id"`
	AssignType   AssignTypeStr `gorm:"not null;unique" json:"assign_type" binding:"enum"`
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

func (b BasicKpiType) IsValid() bool {
	switch b {
	case constants.SINGLE_KPI_TYPE, constants.MULTI_KPI_TYPE:
		return true
	}

	return false
}

func (a AssignTypeStr) IsValid() bool {
	switch a {
	case constants.ASSIGN_TYPE_ROLE, constants.ASSIGN_TYPE_TEAM, constants.ASSIGN_TYPE_INDIVIDUAL:
		return true
	}

	return false
}
