package models

type Appraisal struct {
	CommonModel
	AppraisalName    string         `gorm:"not null" json:"appraisal_name" binding:"required,min=3,max=30"`
	AppraisalYear    uint32         `gorm:"not null" json:"appraisal_year" binding:"required,gte=2023"`
	TeamId           uint64         `gorm:"not null" json:"team_id" binding:"required"`
	AppraisalFlowID  uint64         `gorm:"not null" json:"appraisal_flow_id" binding:"required"`
	AppraisalFlow    AppraisalFlow  `json:"-"`
	SupervisorID     uint64         `gorm:"not null" json:"supervisor_id" binding:"required"`
	AppraisalTypeStr string         `gorm:"not null" json:"appraisal_type" binding:"required"`
	AppraisalType    AppraisalType  `gorm:"references:AppraisalType;foreignKey:AppraisalTypeStr" json:"-"`
	AppraisalKpis    []AppraisalKpi `gorm:"foreignKey:AppraisalID;not null" json:"appraisal_kpis" binding:"required"`
}

type AppraisalKpi struct {
	CommonModel
	AppraisalID uint64 `gorm:"not null" json:"-"`
	EmployeeID  uint64 `gorm:"not null" json:"employee_id" binding:"required"`
	KpiID       uint64 `gorm:"not null" json:"kpi_id" binding:"required"`
	Kpi         Kpi    `json:"-"`
	Status      string `gorm:"not null" json:"status" binding:"required"`
}

type AppraisalType struct {
	CommonModel
	AppraisalType string `gorm:"not null;unique" json:"appraisal_type"`
}
