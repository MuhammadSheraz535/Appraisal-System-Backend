package models

type Appraisal struct {
	CommonModel
	AppraisalName      string         `gorm:"not null" json:"appraisal_name" binding:"required,min=3,max=20"`
	AppraisalYear      uint16         `gorm:"not null" json:"appraisal_year" binding:"required,gte=2023"`
	AppraisalTypeStr   string         `gorm:"not null" json:"appraisal_type" binding:"required"`
	AppraisalType      AppraisalType  `gorm:"references:AppraisalType;foreignKey:AppraisalTypeStr" json:"-"`
	SupervisorID       uint16         `gorm:"not null" json:"supervisor_id" binding:"required"`
	SupervisorName     string         `json:"supervisor_name,omitempty"`
	AppraisalFlowID    uint16         `gorm:"not null" json:"appraisal_flow_id" binding:"required"`
	AppraisalFlow      AppraisalFlow  `json:"appraisal_flow"`
	AppraisalFor       uint16         `gorm:"not null" json:"appraisal_for_id"  binding:"required"`
	AppraisalForName   string         `json:"appraisal_for,omitempty"`
	SelectedFieldID    uint16         `gorm:"not null" json:"selected_id" binding:"required"`
	SelectedFieldNames string         `json:"appraisal_for_name,omitempty"`
	AssignType         AssignType     `gorm:"references:AssignTypeId;foreignKey:AppraisalFor" json:"-"`
	Status             *bool          `gorm:"not null" json:"status" binding:"required"`
	AppraisalKpis      []AppraisalKpi `gorm:"foreignKey:AppraisalID;not null" json:"appraisal_kpis" binding:"required"`
}

type AppraisalKpi struct {
	CommonModel
	AppraisalID uint16 `gorm:"not null" json:"-"`
	EmployeeID  uint16 `gorm:"not null" json:"employee_id" binding:"required"`
	KpiID       uint16 `gorm:"not null" json:"kpi_id" binding:"required"`
	Kpi         Kpi    `json:"-"`
	Status      string `gorm:"not null" json:"status" binding:"required"`
}

type AppraisalType struct {
	CommonModel
	AppraisalType string `gorm:"not null;unique" json:"appraisal_type"`
}
