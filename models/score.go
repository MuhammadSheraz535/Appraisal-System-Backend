package models

type Score struct {
	CommonModel
	AppraisalKpiID uint16       `gorm:"not null" json:"appraisal_kpi_id" validate:"required"`
	AppraisalKpi   AppraisalKpi `json:"appraisal_kpi"`
	EvaluatorID    uint16       `gorm:"not null" json:"evaluator_id"`
	Score          uint16       `json:"score,omitempty"`
	TextAnswer     string       `json:"text_answer,omitempty"`
}
