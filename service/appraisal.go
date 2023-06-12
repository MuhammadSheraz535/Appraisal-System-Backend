package service

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mrehanabbasi/appraisal-system-backend/constants"
	"github.com/mrehanabbasi/appraisal-system-backend/controller"
	"github.com/mrehanabbasi/appraisal-system-backend/database"
	log "github.com/mrehanabbasi/appraisal-system-backend/logger"
	"github.com/mrehanabbasi/appraisal-system-backend/models"
	"github.com/mrehanabbasi/appraisal-system-backend/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AppraisalService struct {
	Db *gorm.DB
}

func NewAppraisalService() *AppraisalService {
	db := database.DB
	err := db.AutoMigrate(&models.Appraisal{}, models.EmployeeData{}, models.AppraisalKpi{}, models.Score{})
	if err != nil {
		panic(err)
	}

	return &AppraisalService{Db: db}
}

func (r *AppraisalService) CreateAppraisal(c *gin.Context) {
	log.Info("Initializing CreateAppraisal handler function...")

	var appraisal models.Appraisal
	err := c.ShouldBindJSON(&appraisal)
	if err != nil {
		errs, ok := controller.ErrValidationSlice(err)
		if !ok {
			log.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Error(err.Error())
		if len(errs) > 1 {
			c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": errs[0]})
		}
		return
	}

	// Validate each FlowStep struct
	for _, ak := range appraisal.AppraisalKpis {
		//check employee id exist in toss api

		if ak.EmployeeID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "employee_id field is required"})
			return
		}

		if ak.KpiID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "kpi_id field is required "})
			return
		}

		if ak.Status == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "status field is required "})
			return
		}
		//check employees id in toss api
		errCode, err := utils.CheckIndividualAgainstToss(ak.EmployeeID)
		if err != nil {
			log.Error(err.Error())
			c.JSON(errCode, gin.H{"error": err.Error()})
			return
		}
	}

	// Fetch all employee IDs from the AppraisalKpis table
	existingEmployeeIDs := make([]int, 0)
	if err := r.Db.Model(&models.AppraisalKpi{}).Pluck("employee_id", &existingEmployeeIDs).Error; err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve existing employee IDs"})
		return
	}
	_, name, err := checkAssignType(r.Db, uint16(appraisal.AppraisalFor))
	if err != nil {
		log.Error("invalid assign type")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assign type"})
		return
	}
	appraisal.AppraisalForName = name

	switch appraisal.AppraisalForName {
	case constants.ASSIGN_TYPE_TEAM:
		errCode, name, err := utils.VerifyTeamAndSupervisorID(appraisal.SelectedFieldID, appraisal.SupervisorID)
		if err != nil {
			log.Error(err.Error())
			c.JSON(errCode, gin.H{"error": err.Error()})
			return
		}
		appraisal.SelectedFieldNames = name
		kpis := make([]models.Kpi, 0)

		empIds, err := utils.GetEmployeesId(uint16(appraisal.SelectedFieldID))
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch employee IDs"})
			return
		}

		// Fetch employee names and role IDs for each employee ID
		employeeDataList := make([]models.EmployeeData, 0)
		for _, empID := range empIds {
			empName, err := utils.GetEmployeeName(empID)
			if err != nil {
				log.Error("Invalid Employee ID")
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Employee ID"})
				return
			}

			roleIDs, err := utils.GetRolesID([]uint16{empID})
			if err != nil {
				log.Error(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch role IDs"})
				return
			}

			roleID := roleIDs[0] // Retrieve the RoleID for the employee (assuming there's only one role ID for each employee)

			designationName, err := utils.GetDesignationName(roleID)
			if err != nil {
				log.Error(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch role IDs"})
				return
			}

			employeeImage, err := utils.GetEmployeeImageByID(uint64(empID))
			if err != nil {
				log.Error(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch employee image"})
				return
			}

			projectDetails, err := utils.GetProjectDetailsByEmployeeID(empID)
			if err != nil {
				log.Error(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch project details"})
				return
			}

			var ProjectID uint16
			var ProjectName string

			for _, project := range projectDetails {
				ProjectID = project.ProjectDetails.ProjectID
				ProjectName = project.ProjectDetails.ProjectName
			}
			baseurl := os.Getenv("TOSS_BASE_URL")
			employeeData := models.EmployeeData{
				AppraisalID:     appraisal.ID,
				TossEmpID:       empID,
				EmployeeName:    empName,
				TeamID:          ProjectID,
				TeamName:        ProjectName,
				EmployeeImage:   baseurl + "/" + employeeImage,
				Designation:     roleID, // Assign the RoleID as Designation
				DesignationName: designationName,
				AppraisalStatus: "pending",
			}
			employeeDataList = append(employeeDataList, employeeData)
		}

		// Append EmployeeData to Appraisal
		appraisal.EmployeesList = employeeDataList

		roleIds, err := utils.GetRolesID(empIds)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch role IDs"})
			return
		}

		db := r.Db.Model(&models.Kpi{})
		db = db.Joins("JOIN assign_types ON assign_types.id = kpis.assign_type_id").
			Where(`(kpis.selected_assign_id = ? AND assign_types.assign_type = ?)
			OR (kpis.selected_assign_id IN (?) AND assign_types.assign_type = ?)
			OR (kpis.selected_assign_id IN (?) AND assign_types.assign_type = ?)`,
				appraisal.SelectedFieldID, constants.ASSIGN_TYPE_TEAM, empIds, constants.ASSIGN_TYPE_INDIVIDUAL, roleIds, constants.ASSIGN_TYPE_ROLE)
		err = db.Model(&models.Kpi{}).Order("id ASC").Find(&kpis).Error
		if err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		if len(kpis) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "KPI does not exist for the team"})
			return
		}

		for _, kpi := range kpis {
			if kpi.AssignTypeName == constants.ASSIGN_TYPE_INDIVIDUAL {
				appraisalKpi := models.AppraisalKpi{
					AppraisalID: appraisal.ID,
					EmployeeID:  kpi.SelectedAssignID,
					KpiID:       kpi.ID,
					Status:      "pending",
				}
				appraisal.AppraisalKpis = append(appraisal.AppraisalKpis, appraisalKpi)
			}

			if kpi.AssignTypeName == constants.ASSIGN_TYPE_TEAM {

				// Assign the KPI to individual employees of the team
				for _, employeeID := range empIds {
					appraisalKpi := models.AppraisalKpi{
						AppraisalID: appraisal.ID,
						EmployeeID:  employeeID,
						KpiID:       kpi.ID,
						Status:      "pending",
					}
					appraisal.AppraisalKpis = append(appraisal.AppraisalKpis, appraisalKpi)
				}
			}
			if kpi.AssignTypeName == constants.ASSIGN_TYPE_ROLE {

				for _, employeeID := range empIds {
					roleIDs, err := utils.GetRolesID([]uint16{employeeID})
					if err != nil {
						log.Error(err)
						c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch role IDs"})
						return
					}
					roleID := roleIDs[0] // Assuming there's only one role ID for each employee

					if roleID == kpi.SelectedAssignID {
						appraisalKpi := models.AppraisalKpi{
							AppraisalID: appraisal.ID,
							EmployeeID:  employeeID,
							KpiID:       kpi.ID,
							Status:      "pending",
						}
						appraisal.AppraisalKpis = append(appraisal.AppraisalKpis, appraisalKpi)
					}
				}

			}

		}

	case constants.ASSIGN_TYPE_INDIVIDUAL:
		errCode, name, err := utils.VerifyIndividualAndSupervisorID(appraisal.SelectedFieldID, appraisal.SupervisorID)
		if err != nil {
			log.Error(err.Error())
			c.JSON(errCode, gin.H{"error": err.Error()})
			return
		}

		// Fetch employee name
		empName, err := utils.GetEmployeeName(appraisal.SelectedFieldID)
		if err != nil {
			log.Error("Invalid Employee ID")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Employee ID"})
			return
		}

		// Get the role ID for the employee
		roleIDs, err := utils.GetRolesID([]uint16{uint16(appraisal.SelectedFieldID)})
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch role IDs"})
			return
		}

		roleID := roleIDs[0] // Retrieve the RoleID for the employee (assuming there's only one role ID for each employee)

		designationName, err := utils.GetDesignationName(uint16(roleID))
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch role IDs"})
			return
		}

		employeeImage, err := utils.GetEmployeeImageByID(uint64(appraisal.SelectedFieldID))
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch Employees Image"})
			return
		}

		projectDetails, err := utils.GetProjectDetailsByEmployeeID(appraisal.SelectedFieldID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch Project Details"})
			log.Error(err.Error())
			return
		}

		var ProjectID uint16
		var ProjectName string

		for _, project := range projectDetails {
			ProjectID = project.ProjectDetails.ProjectID
			ProjectName = project.ProjectDetails.ProjectName
		}
		baseurl := os.Getenv("TOSS_BASE_URL")
		// Create EmployeeData instance
		employeeData := models.EmployeeData{
			AppraisalID:     appraisal.ID,
			TossEmpID:       appraisal.AppraisalFor,
			EmployeeName:    empName,
			TeamID:          ProjectID,
			TeamName:        ProjectName,
			EmployeeImage:   baseurl + "/" + employeeImage,
			Designation:     roleID, // Assign the RoleID as Designation
			DesignationName: designationName,
			AppraisalStatus: "pending",
		}

		// Append EmployeeData to Appraisal
		appraisal.EmployeesList = []models.EmployeeData{employeeData}

		appraisal.SelectedFieldNames = name
		kpis := make([]models.Kpi, 0)
		if err := r.Db.Where("assign_type_id = ? AND selected_assign_id = ?", appraisal.AppraisalFor, appraisal.SelectedFieldID).Find(&kpis).Error; err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while retrieving KPIs"})
			return
		}

		if len(kpis) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Kpi does not exist for the Individual"})
			return
		}

		for _, kpi := range kpis {
			appraisalKpi := models.AppraisalKpi{
				AppraisalID: appraisal.ID,
				EmployeeID:  appraisal.SelectedFieldID,
				KpiID:       kpi.ID,
				Status:      "pending",
			}
			appraisal.AppraisalKpis = append(appraisal.AppraisalKpis, appraisalKpi)
		}

		// Modify the Case ROLE section
	case constants.ASSIGN_TYPE_ROLE:
		errCode, name, err := utils.CheckRoleExists(appraisal.SelectedFieldID)
		if err != nil {
			log.Error(err.Error())
			c.JSON(errCode, gin.H{"error": err.Error()})
			return
		}
		appraisal.SelectedFieldNames = name

		// Get employee IDs for the provided role ID
		employeeIDs, err := utils.GetEmployeeIDsByDesignation(uint16(appraisal.SelectedFieldID))
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch employee IDs"})
			return
		}

		if len(employeeIDs) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "No employees found for the provided role"})
			return
		}

		employeeDataList := make([]models.EmployeeData, 0)
		for _, empID := range employeeIDs {
			empName, err := utils.GetEmployeeName(empID)
			if err != nil {
				log.Error("Invalid Employee ID")
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Employee ID"})
				return
			}

			designationName, err := utils.GetDesignationName(uint16(appraisal.SelectedFieldID))
			if err != nil {
				log.Error(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch designation name"})
				return
			}
			employeeImage, err := utils.GetEmployeeImageByID(uint64(empID))
			if err != nil {
				log.Error(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch Employee Image"})
				return
			}

			projectDetails, err := utils.GetProjectDetailsByEmployeeID(empID)
			if err != nil {
				log.Error(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch Project Details"})
				return
			}

			var ProjectID uint16
			var ProjectName string

			for _, project := range projectDetails {
				ProjectID = project.ProjectDetails.ProjectID
				ProjectName = project.ProjectDetails.ProjectName
			}
			baseurl := os.Getenv("TOSS_BASE_URL")
			employeeData := models.EmployeeData{
				AppraisalID:     appraisal.ID,
				TossEmpID:       empID,
				EmployeeName:    empName,
				EmployeeImage:   baseurl + "/" + employeeImage,
				Designation:     uint16(appraisal.SelectedFieldID),
				DesignationName: designationName,
				TeamID:          ProjectID,
				TeamName:        ProjectName,
				AppraisalStatus: "pending",
			}
			employeeDataList = append(employeeDataList, employeeData)
		}

		// Append EmployeeData to Appraisal
		appraisal.EmployeesList = employeeDataList

		kpis := make([]models.Kpi, 0)
		if err := r.Db.Where("assign_type_id = ? AND selected_assign_id = ?", appraisal.AppraisalFor, appraisal.SelectedFieldID).Find(&kpis).Error; err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while retrieving KPIs"})
			return
		}

		if len(kpis) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "KPI does not exist for the Role"})
			return
		}

		for _, kpi := range kpis {
			for _, empID := range employeeIDs {
				appraisalKpi := models.AppraisalKpi{
					AppraisalID: appraisal.ID,
					EmployeeID:  empID,
					KpiID:       kpi.ID,
					Status:      "pending",
				}
				appraisal.AppraisalKpis = append(appraisal.AppraisalKpis, appraisalKpi)
			}
		}
	}

	// check appraisal type exists
	err = checkAppraisalType(r.Db, appraisal.AppraisalTypeStr)
	if err != nil {
		log.Error("invalid appraisal type")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid appraisal type"})
		return
	}
	var appraisalFlow models.AppraisalFlow
	err = r.Db.Model(&models.AppraisalFlow{}).First(&appraisalFlow, appraisal.AppraisalFlowID).Error
	if err != nil {
		log.Error("invalid appraisal flow ID")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid appraisal flow ID"})
		return
	}

	// Call GetSupervisorName function to retrieve the supervisor name
	supervisorName, err := utils.GetSupervisorName(appraisal.SupervisorID)
	if err != nil {
		log.Error("failed to get supervisor name")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get supervisor name"})
		return
	}
	appraisal.SupervisorName = supervisorName

	dbAppraisal, err := controller.CreateAppraisal(r.Db, &appraisal)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dbAppraisal)
}

func (r *AppraisalService) GetAppraisalByID(c *gin.Context) {
	log.Info("Initializing GetAppraisalByID handler function...")

	id, _ := strconv.ParseUint(c.Param("id"), 0, 64)

	var appraisal models.Appraisal
	err := controller.GetAppraisalByID(r.Db, &appraisal, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})

		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		return
	}

	// Create a response structure with the required fields
	response := struct {
		Appraisal models.Appraisal `json:"appraisal"`
	}{
		Appraisal: appraisal,
	}

	c.JSON(http.StatusOK, response)
}

func (r *AppraisalService) GetEmployeeDataByAppraisalID(c *gin.Context) {
	log.Info("Initializing GetEmployeeDataByAppraisalID handler function...")

	id, _ := strconv.ParseUint(c.Param("id"), 0, 64)
	var employeeData []models.EmployeeData
	db := r.Db.Model(&models.EmployeeData{})

	tossEmpId := c.Query("toss_emp_id")
	if tossEmpId != "" {
		db = db.Where("toss_emp_id = ?", tossEmpId)
	}

	err := controller.GetEmployeeDataByAppraisalID(db, &employeeData, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})

		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		return
	}

	c.JSON(http.StatusOK, employeeData)
}

func (r *AppraisalService) GetAllAppraisals(c *gin.Context) {
	log.Info("Initializing GetAllAppraisal handler function...")
	var appraisals []models.Appraisal
	db := r.Db.Model(&models.Appraisal{})

	appraisalName := c.Query("appraisal_name")
	supervisorID := c.Query("supervisor_id")

	if appraisalName != "" {
		db = db.Where("appraisal_name LIKE ?", "%"+appraisalName+"%")
	}

	if supervisorID != "" {
		db = db.Where("supervisor_id = ?", supervisorID)
	}

	err := controller.GetAllAppraisals(db, &appraisals)

	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appraisals)
}

func (r *AppraisalService) UpdateAppraisal(c *gin.Context) {
	log.Info("Initializing UpdateAppraisal handler function...")

	id, _ := strconv.ParseUint(c.Param("id"), 0, 16)

	var appraisal models.Appraisal

	err := c.ShouldBindJSON(&appraisal)
	if err != nil {
		errs, ok := controller.ErrValidationSlice(err)
		if !ok {
			log.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Error(err.Error())
		if len(errs) > 1 {
			c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": errs[0]})
		}
		return
	}

	// Validate each FlowStep struct
	for _, ak := range appraisal.AppraisalKpis {
		if ak.EmployeeID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "employee_id field is required"})
			return
		}

		if ak.KpiID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "kpi_id field is required "})
			return
		}

		if ak.Status == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "status field is required "})
			return
		}
		//check employees id in toss api
		errCode, err := utils.CheckIndividualAgainstToss(ak.EmployeeID)
		if err != nil {
			log.Error(err.Error())
			c.JSON(errCode, gin.H{"error": err.Error()})
			return
		}
	}

	// Fetch all employee IDs from the AppraisalKpis table
	existingEmployeeIDs := make([]int, 0)
	if err := r.Db.Model(&models.AppraisalKpi{}).Pluck("employee_id", &existingEmployeeIDs).Error; err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve existing employee IDs"})
		return
	}

	// Check if the provided employee IDs exist in the AppraisalKpis table
	for _, ed := range appraisal.EmployeesList {

		errCode, _, err := utils.CheckRoleExists(ed.Designation)
		if err != nil {
			log.Error(err.Error())
			c.JSON(errCode, gin.H{"error": err.Error()})
			return
		}

		employeeIDExists := false
		for _, existingID := range existingEmployeeIDs {
			if existingID == int(ed.TossEmpID) {
				employeeIDExists = true
				break
			}
		}

		if !employeeIDExists {
			c.JSON(http.StatusNotFound, gin.H{"error": "No KPI found against this Employee"})
			return
		}

		// Check employee ID in the Toss API
		errCode, err = utils.CheckIndividualAgainstToss(ed.TossEmpID)
		if err != nil {
			log.Error(err.Error())
			c.JSON(errCode, gin.H{"error": err.Error()})
			return
		}
	}

	//check assigns type
	_, name, err := checkAssignType(r.Db, uint16(appraisal.AppraisalFor))

	if err != nil {
		log.Error("invalid assign type")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid assign type"})
		return
	}
	appraisal.AppraisalForName = name

	switch appraisal.AppraisalForName {
	case constants.ASSIGN_TYPE_TEAM:
		errCode, name, err := utils.VerifyTeamAndSupervisorID(appraisal.SelectedFieldID, appraisal.SupervisorID)
		if err != nil {
			log.Error(err.Error())
			c.JSON(errCode, gin.H{"error": err.Error()})
			return
		}
		appraisal.SelectedFieldNames = name

	case constants.ASSIGN_TYPE_INDIVIDUAL:
		errCode, name, err := utils.VerifyIndividualAndSupervisorID(appraisal.SelectedFieldID, appraisal.SupervisorID)
		if err != nil {
			log.Error(err.Error())
			c.JSON(errCode, gin.H{"error": err.Error()})
			return
		}
		appraisal.SelectedFieldNames = name

	case constants.ASSIGN_TYPE_ROLE:
		errCode, name, err := utils.CheckRoleExists(appraisal.SelectedFieldID)
		if err != nil {
			log.Error(err.Error())
			c.JSON(errCode, gin.H{"error": err.Error()})
			return
		}
		appraisal.SelectedFieldNames = name
	}

	err = checkAppraisalType(r.Db, appraisal.AppraisalTypeStr)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid appraisal type"})
		return
	}
	appraisal.ID = uint16(id)

	// checking appraisal flow id exists in db
	var appraisalFlow models.AppraisalFlow
	err = r.Db.Model(&models.AppraisalFlow{}).First(&appraisalFlow, appraisal.AppraisalFlowID).Error
	if err != nil {
		log.Error("invalid appraisal flow id")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid appraisal flow id"})
		return
	}

	// Call GetSupervisorName function to retrieve the supervisor name
	supervisorName, err := utils.GetSupervisorName(appraisal.SupervisorID)
	if err != nil {
		log.Error("failed to get supervisor name")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get supervisor name"})
		return
	}
	appraisal.SupervisorName = supervisorName

	// callling controller update function
	dbAppraisal, err := controller.UpdateAppraisal(r.Db, &appraisal)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dbAppraisal)
}

func (r *AppraisalService) DeleteAppraisal(c *gin.Context) {
	log.Info("Initializing DeleteAppraisal handler function...")

	var appraisal models.Appraisal
	id, _ := strconv.ParseUint(c.Param("id"), 0, 16)
	appraisal.ID = uint16(id)

	err := controller.GetAppraisalByID(r.Db, &appraisal, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})

		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		return
	}

	err = controller.DeleteAppraisal(r.Db, &appraisal, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func checkAppraisalType(db *gorm.DB, appraisal_type string) error {
	log.Info("Checking Appraisal type")
	var appraisalTypeModel models.AppraisalType
	err := db.Model(&models.AppraisalType{}).Where("appraisal_type = ?", appraisal_type).First(&appraisalTypeModel).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *AppraisalService) GetAppraisalKpisByEmpID(c *gin.Context) {
	log.Info("Initializing GetAppraisalkpisByEmpID handler function...")

	id, _ := strconv.ParseUint(c.Param("emp_id"), 0, 64)
	var appraisalKpi []models.AppraisalKpi

	//Adding query parameters for employees id
	db := r.Db.Model(&models.AppraisalKpi{})
	employeeid := c.Query("employee_id")
	if employeeid != "" {
		db = db.Where("employee_id = ?", employeeid)
	}

	err := controller.GetAppraisalKpisByEmpID(db, &appraisalKpi, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})

		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		return
	}

	c.JSON(http.StatusOK, appraisalKpi)
}

func (r *AppraisalService) Score(c *gin.Context) {
	log.Info("Initializing Score handler function...")

	appraisalID := c.Param("id")
	employeeID := c.Param("emp_id")

	// Check if the employee ID exists
	var appraisalKpi models.AppraisalKpi
	err := r.Db.Model(&models.AppraisalKpi{}).Where("appraisal_id = ? AND employee_id = ?", appraisalID, employeeID).First(&appraisalKpi).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Appraiasl not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	var score models.Score

	err = c.ShouldBindJSON(&score)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the appraisal KPI data from the database
	err = r.Db.Model(&models.AppraisalKpi{}).Preload(clause.Associations).Where("appraisal_id = ? AND employee_id = ?", appraisalID, employeeID).First(&appraisalKpi).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Check KPI type and set values accordingly
	switch appraisalKpi.Kpi.KpiTypeStr {
	case constants.FEEDBACK_KPI_TYPE, constants.OBSERVATORY_KPI_TYPE:
		fmt.Println("KPI Type :", appraisalKpi.Kpi.KpiTypeStr)
		score.Score = 0

	case constants.QUESTIONNAIRE_KPI_TYPE:
		fmt.Println("KPI Type: ", appraisalKpi.Kpi.KpiTypeStr)
		if score.Score != 0 && score.Score != 1 {
			errMsg := "Questionnaire Score should be either 0 or 1"
			log.Error(errMsg)
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
			return
		}
		score.TextAnswer = ""

	case constants.MEASURED_KPI_TYPE:
		fmt.Println("KPI Type: ", appraisalKpi.Kpi.KpiTypeStr)
		score.TextAnswer = ""
	}

	// Save the score to the database or perform any necessary operations

	scores, err := controller.Score(r.Db, &score)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, scores)
}
