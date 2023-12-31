package models

import (
	"reflect"
	"strings"

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
	AssignTypeId uint16        `gorm:"not null;unique" json:"assign_type_id"`
	AssignType   AssignTypeStr `gorm:"not null;unique" json:"assign_type"`
}

type Kpi struct {
	CommonModel
	KpiName            string                  `gorm:"size:100;not null;default:''" json:"kpi_name" validate:"required,min=3,max=30"`
	KpiDescription     string                  `gorm:"not null;default:''" json:"kpi_description" validate:"required"`
	AssignTypeID       uint16                  `gorm:"not null;default:0" json:"assign_type_id" validate:"required"`
	AssignTypeName     string                  `json:"assign_type_name,omitempty"`
	SelectedAssignID   uint16                  `gorm:"not null;default:0" json:"selected_assign_id" validate:"required"`
	SelectedAssignName string                  `json:"selected_assign_name,omitempty"`
	AssignType         AssignType              `gorm:"references:AssignTypeId;foreignKey:AssignTypeID" json:"-"`
	KpiTypeStr         string                  `gorm:"not null;default:''" json:"kpi_type" validate:"required"`
	KpiType            KpiType                 `gorm:"references:KpiType;foreignKey:KpiTypeStr" json:"-"`
	KpiWeight          uint8                   `gorm:"not null;default:0" json:"kpi_weight" validate:"required"`
	ApplicableFor      pq.StringArray          `gorm:"type:text[];not null" json:"applicable_for" validate:"required"`
	Statement          string                  `json:"statement,omitempty"`
	Statements         []MultiStatementKpiData `gorm:"foreignKey:KpiID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"statements,omitempty"`
}

type MultiStatementKpiData struct {
	CommonModel
	KpiID     uint16 `gorm:"not null;default:0" json:"-"`
	Statement string `gorm:"not null;default:''" json:"statement" validate:"required"`
	Weightage uint8  `gorm:"not null;default:0" json:"weightage" validate:"required"`
	// CorrectAnswer string `gorm:"not null" json:"correct_answer" validate:"required"`
}

func (a *Kpi) Validate() error {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return validate.Struct(a)
}

func (a *MultiStatementKpiData) Validate() error {
	validate := validator.New()
	return validate.Struct(a)
}
