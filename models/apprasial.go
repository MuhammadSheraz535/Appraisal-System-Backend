package models

type Apprasial struct {
	CommonModel
	ApprasialName   string `json:"apprasial_name" binding:"required" gorm:"not null"`
	TeamId          uint64 `json:"team_id" binding:"required"`
	AppraisalFlowID uint64 `json:"appraisal_flow_id" binding:"required" gorm:"not null"`
	AppraisalFlow   AppraisalFlow
	SupervisorID    uint64 `json:"supervisor_id" binding:"required"`
	AppraisalType   string `json:"appraisal_type" binding:"required" gorm:"not null"`
	AppraisalKpis   []AppraisalKpis
}

type AppraisalKpis struct {
	CommonModel
	EmployeeID  uint64 `json:"employee_id" binding:"required"`
	KpiID       uint64 `json:"kpi_id" binding:"required"`
	Kpi         Kpi
	ApprasialId uint64 `json:"appraisal_id" binding:"required"`
	Status      string `json:"status" binding:"required"`
}
