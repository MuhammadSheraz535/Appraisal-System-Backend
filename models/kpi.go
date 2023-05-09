package models

import (
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
)

type BasicKpiType string
type AssignTypeStr string

type KpiType struct {
	CommonModel
	KpiType      string       `gorm:"not null;unique" json:"kpi_type"`
	BasicKpiType BasicKpiType `gorm:"not null" json:"basic_kpi_type"`
}

type AssignType struct {
	CommonModel
	AssignTypeId uint64        `gorm:"not null;unique" json:"assign_type_id"`
	AssignType   AssignTypeStr `gorm:"not null;unique" json:"assign_type"`
}

type Kpi struct {
	CommonModel
	KpiName        string                  `gorm:"size:100;not null" json:"kpi_name" validate:"required,min=3,max=30"`
	KpiDescription string                  `gorm:"not null" json:"kpi_description" validate:"required"`
	AssignTypeID   uint64                  `gorm:"not null" json:"assign_type" validate:"required"`
	AssignType     AssignType              `gorm:"references:AssignTypeId;foreignKey:AssignTypeID" json:"-"`
	KpiTypeStr     string                  `gorm:"not null" json:"kpi_type" validate:"required"`
	KpiType        KpiType                 `gorm:"references:KpiType;foreignKey:KpiTypeStr" json:"-"`
	ApplicableFor  pq.StringArray          `gorm:"type:text[];not null" json:"applicable_for" validate:"required"`
	Statement      string                  `json:"statement,omitempty"`
	Statements     []MultiStatementKpiData `gorm:"foreignKey:KpiID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"statements,omitempty"`
}

type MultiStatementKpiData struct {
	CommonModel
	KpiID         uint64 `gorm:"not null" json:"-"`
	Statement     string `gorm:"not null" json:"statement" validate:"required"`
	CorrectAnswer string `gorm:"not null" json:"correct_answer" validate:"required"`
	Weightage     uint64 `gorm:"not null" json:"weightage" validate:"required"`
}

func (a *Kpi) Validate() error {
	validate := validator.New()
	return validate.Struct(a)
}

func (a *MultiStatementKpiData) Validate() error {
	validate := validator.New()
	return validate.Struct(a)
}
