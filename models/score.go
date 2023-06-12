package models

type Score struct {
	CommonModel
	AppraisalKpiID uint16       `gorm:"not null;default:0" json:"appraisal_kpi_id" validate:"required"`
	AppraisalKpi   AppraisalKpi `json:"appraisal_kpi"`
	EvaluatorID    uint16       `gorm:"not null;default:0" json:"evaluator_id"`
	Score          *uint16      `gorm:"default:0" json:"score,omitempty"`
	TextAnswer     string       `gorm:";default:''"  json:"text_answer,omitempty"`
}
