package models

type Apprasial struct {
	CommonModel
	ApprasialName    string          `json:"apprasial_name" binding:"required" gorm:"not null"`
	TeamId           uint64          `json:"team_id" binding:"required"`
	AppraisalFlowID  uint64          `json:"appraisal_flow_id" binding:"required" gorm:"not null"`
	AppraisalFlow    AppraisalFlow   `json:"-"`
	SupervisorID     uint64          `json:"supervisor_id" binding:"required"`
	AppraisalTypeStr string          `gorm:"not null" json:"apprasial_type"`
	AppraisalType    AppraisalType   `gorm:"references:AppraisalType;foreignKey:AppraisalTypeStr" json:"-"`
	AppraisalKpis    []AppraisalKpis `json:"appraisal_kpis" gorm:"foreignKey:ApprasialID;not null"`
}

type AppraisalKpis struct {
	CommonModel
	ApprasialID uint64 `json:"-"`
	EmployeeID  uint64 `json:"employee_id"`
	KpiID       uint64 `json:"kpi_id" gorm:"not null"`
	Kpi         Kpi    `json:"-"`
	Status      string `json:"status"`
}

type AppraisalType struct {
	CommonModel
	AppraisalType string `gorm:"not null;unique" json:"appraisal_type"`
}
