package models

type Score struct {
	CommonModel
	AppraisalKpiID uint64       `gorm:"not null" json:"appraisal_kpi_id"`
	AppraisalKpi   AppraisalKpi `json:"-"`
	EvaluatorID    uint64       `gorm:"not null" json:"evaluator_id"`
	Score          uint32       `gorm:"not null" json:"score"`
	TextAnswer     string       `gorm:"not null" json:"text_answer"`
}
