package models

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type AppraisalFlow struct {
	CommonModel
	FlowName         string        `json:"flow_name" gorm:"type:varchar(255);not null;unique" validate:"required,min=3,max=30"`
	CreatedBy        uint64        `json:"created_by" gorm:"not null" validate:"required"`
	IsActive         bool          `json:"is_active" gorm:"not null" validate:"required"`
	TeamId           uint64        `json:"team_id" gorm:"not null" validate:"required"`
	AppraisalTypeStr string        `json:"appraisal_type" gorm:"not null" binding:"required"`
	AppraisalType    AppraisalType `json:"-" gorm:"references:AppraisalType;foreignKey:AppraisalTypeStr"`
	FlowSteps        []FlowStep    `json:"flow_steps" gorm:"foreignKey:FlowID;not null" validate:"required"`
}

type FlowStep struct {
	CommonModel
	FlowID    uint64 `json:"-" gorm:"not null"`
	StepName  string `json:"step_name"  gorm:"type:varchar(255);not null" validate:"required"`
	StepOrder uint64 `json:"step_order"  gorm:"not null" validate:"required"`
	UserId    uint64 `json:"user_id"  gorm:"not null" validate:"required"`
}

func (a *AppraisalFlow) Validate() error {
	validate := validator.New()
	return validate.Struct(a)
}

func (f *FlowStep) Validate() error {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return validate.Struct(f)
}
