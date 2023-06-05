package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/mrehanabbasi/appraisal-system-backend/constants"
	log "github.com/mrehanabbasi/appraisal-system-backend/logger"
)

func VerifyIdAgainstTossApis(selectedAssignID uint16, assignType string) (int, string, error) {
	// Check which SelectedAssignID exists in the API
	tossBaseUrl := os.Getenv("TOSS_BASE_URL")

	var selectedAssignName string
	switch assignType {
	case constants.ASSIGN_TYPE_ROLE:
		method := http.MethodGet
		url := tossBaseUrl + "/api/Employee/GetDesignationsList"

		resp, err := SendRequest(method, url, nil)
		if err != nil {
			log.Error(err.Error())
			return http.StatusInternalServerError, "", err
		}
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error(err.Error())
			return http.StatusInternalServerError, "", err
		}

		type SystemRole struct {
			Value uint16 `json:"value"`
			Label string `json:"label"`
		}

		type DesignationRolesResponse struct {
			SystemRoles []SystemRole `json:"designations"`
		}

		var response DesignationRolesResponse
		if err := json.Unmarshal(responseBody, &response); err != nil {
			log.Error(err.Error())
			return http.StatusInternalServerError, "", err
		}

		found := false
		for _, role := range response.SystemRoles {
			if role.Value == selectedAssignID {
				found = true
				selectedAssignName = role.Label
				break
			}
		}

		if !found {
			err := errors.New("invalid selected role id")
			log.Error(err.Error())
			return http.StatusBadRequest, selectedAssignName, err
		}
	case constants.ASSIGN_TYPE_TEAM:
		method := http.MethodGet
		url := tossBaseUrl + "/api/Project/GetAllProjects"

		resp, err := SendRequest(method, url, nil)
		if err != nil {
			log.Error(err.Error())
			return http.StatusInternalServerError, "", err
		}
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error(err.Error())
			return http.StatusInternalServerError, "", err
		}

		var projects []struct {
			ProjectDetails struct {
				ProjectID uint16 `json:"projectId"`
				TeamName  string `json:"projectName"`
			} `json:"projectDetails"`
		}

		if err := json.Unmarshal(responseBody, &projects); err != nil {
			log.Error(err.Error())
			return http.StatusInternalServerError, "", err
		}

		found := false
		for _, project := range projects {
			if project.ProjectDetails.ProjectID == selectedAssignID {
				found = true
				selectedAssignName = project.ProjectDetails.TeamName
				break
			}
		}

		if !found {
			err := errors.New("invalid selected team id")
			log.Error(err.Error())
			return http.StatusBadRequest, selectedAssignName, err
		}
	case constants.ASSIGN_TYPE_INDIVIDUAL:
		method := http.MethodGet
		url := tossBaseUrl + "/api/Employee/GetAllEmployees?AllEmployees=true"

		resp, err := SendRequest(method, url, nil)
		if err != nil {
			log.Error(err.Error())
			return http.StatusInternalServerError, "", err
		}
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error(err.Error())
			return http.StatusInternalServerError, "", err
		}

		var Employees []struct {
			EmployeeID   uint16 `json:"employeeId"`
			EmployeeName string `json:"name"`
		}
		if err := json.Unmarshal(responseBody, &Employees); err != nil {
			log.Error(err.Error())
			return http.StatusInternalServerError, "", err
		}

		found := false
		for _, employee := range Employees {
			if employee.EmployeeID == selectedAssignID {
				found = true
				selectedAssignName = employee.EmployeeName
				break
			}
		}

		if !found {
			err := errors.New("invalid selected employee id")
			log.Error(err.Error())
			return http.StatusBadRequest, "", err
		}
	}

	return 0, selectedAssignName, nil
}

func CheckIndividualAgainstToss(CreatedBy uint16) (int, error) {

	tossBaseUrl := os.Getenv("TOSS_BASE_URL")

	method := http.MethodGet
	url := tossBaseUrl + "/api/Employee/GetAllEmployees?AllEmployees=true"
	resp, err := SendRequest(method, url, nil)
	if err != nil {
		log.Error(err.Error())
		return http.StatusInternalServerError, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err.Error())
		return http.StatusInternalServerError, err
	}

	var Employees []struct {
		EmployeeID uint16 `json:"employeeId"`
	}
	if err := json.Unmarshal(responseBody, &Employees); err != nil {
		log.Error(err.Error())
		return http.StatusInternalServerError, err
	}

	found := false
	for _, employee := range Employees {
		if employee.EmployeeID == CreatedBy {
			found = true
			break
		}
	}

	if !found {
		err := errors.New("invalid selected employee id")
		log.Error(err.Error())
		return http.StatusBadRequest, err
	}

	return 0, nil
}

func GetSupervisorName(SprID uint16) (string, error) {
	tossBaseUrl := os.Getenv("TOSS_BASE_URL") // Get the TOSS base URL from the environment variable
	method := http.MethodGet                  // HTTP method for sending the request

	url := tossBaseUrl + "/api/Project/GetAllProjects" // Construct the URL for fetching all projects

	resp, err := SendRequest(method, url, nil) // Send the HTTP request to the specified URL
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg := "Failed to get supervisor name for supervisor ID: " + strconv.Itoa(int(SprID)) + ". Status code: " + strconv.Itoa(resp.StatusCode)
		return "", errors.New(errMsg) // Return an error if the response status code is not OK
	}

	var projects []struct {
		ProjectDetails struct {
			SupervisorID   uint16 `json:"supervisorId"`
			SupervisorName string `json:"supervisorName"`
		} `json:"projectDetails"`
	}

	responseBody, err := io.ReadAll(resp.Body) // Read the response body
	if err != nil {
		return "", err // Return an error if there's an issue reading the response body
	}

	if err := json.Unmarshal(responseBody, &projects); err != nil {
		return "", err // Return an error if there's an issue unmarshaling the JSON response
	}

	for _, project := range projects {
		if project.ProjectDetails.SupervisorID == SprID {
			return project.ProjectDetails.SupervisorName, nil // Return the Supervisor Name if the supervisor ID matches
		}
	}

	return "", errors.New("supervisor not found") // Return an error if the supervisor ID is not found in the projects
}

func CheckRoleExists(AppraisalForID uint16) (int, string, error) {
	tossBaseUrl := os.Getenv("TOSS_BASE_URL")
	method := http.MethodGet
	url := tossBaseUrl + "/api/Employee/GetDesignationsList"

	resp, err := SendRequest(method, url, nil)
	if err != nil {
		log.Error(err.Error())
		return http.StatusInternalServerError, "", err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err.Error())
		return http.StatusInternalServerError, "", err
	}

	type SystemRole struct {
		Value uint16 `json:"value"`
		Label string `json:"label"`
	}
	var roleName string

	type DesignationRoleResponse struct {
		SystemRoles []SystemRole `json:"designations"`
	}

	var response DesignationRoleResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		log.Error(err.Error())
		return http.StatusInternalServerError, "", err
	}

	found := false
	for _, role := range response.SystemRoles {
		if role.Value == AppraisalForID {
			found = true
			roleName = role.Label
			break
		}
	}

	if !found {
		err := errors.New("invalid selected role id")
		log.Error(err.Error())
		return http.StatusBadRequest, "", err
	}

	return 0, roleName, nil
}

func GetEmployeeName(employeeID uint16) (string, error) {
	tossBaseUrl := os.Getenv("TOSS_BASE_URL") // Get the TOSS base URL from the environment variable
	method := http.MethodGet                  // HTTP method for sending the request

	url := tossBaseUrl + "/api/Employee/GetAllEmployees?AllEmployees=true" // Construct the URL for fetching all employees

	resp, err := SendRequest(method, url, nil) // Send the HTTP request to the specified URL
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg := "Failed to get employee name for employee ID: " + strconv.Itoa(int(employeeID)) + ". Status code: " + strconv.Itoa(resp.StatusCode)
		return "", errors.New(errMsg) // Return an error if the response status code is not OK
	}

	var employees []struct {
		EmployeeID uint16 `json:"employeeId"`
		Name       string `json:"name"`
	}

	responseBody, err := io.ReadAll(resp.Body) // Read the response body
	if err != nil {
		return "", err // Return an error if there's an issue reading the response body
	}

	if err := json.Unmarshal(responseBody, &employees); err != nil {
		return "", err // Return an error if there's an issue unmarshaling the JSON response
	}

	for _, emp := range employees {
		if emp.EmployeeID == employeeID {
			fmt.Printf("Role name is %v", emp.Name)
			return emp.Name, nil // Return the Employee Name if the employee ID matches

		}

	}

	return "", errors.New("employee not found") // Return an error if the employee ID is not found in the employees
}

func GetRoleName(roleID uint16) (string, error) {
	tossBaseUrl := os.Getenv("TOSS_BASE_URL") // Get the TOSS base URL from the environment variable
	method := http.MethodGet                  // HTTP method for sending the request

	url := tossBaseUrl + "/api/Employee/GetSystemRolesList" // Construct the URL for fetching system roles

	resp, err := SendRequest(method, url, nil) // Send the HTTP request to the specified URL
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg := "Failed to get role name for role ID: " + strconv.Itoa(int(roleID)) + ". Status code: " + strconv.Itoa(resp.StatusCode)
		return "", errors.New(errMsg) // Return an error if the response status code is not OK
	}

	type SystemRole struct {
		Value uint16 `json:"value"`
		Label string `json:"label"`
	}

	type SystemRolesResponse struct {
		SystemRoles []SystemRole `json:"systemRoles"`
	}

	var response SystemRolesResponse
	responseBody, err := io.ReadAll(resp.Body) // Read the response body
	if err != nil {
		return "", err // Return an error if there's an issue reading the response body
	}

	if err := json.Unmarshal(responseBody, &response); err != nil {
		return "", err // Return an error if there's an issue unmarshaling the JSON response
	}

	for _, role := range response.SystemRoles {
		if role.Value == roleID {
			fmt.Printf("Role name is %v", role.Label)
			return role.Label, nil // Return the Role Name if the role ID matches
		}

	}

	return "", errors.New("role not found") // Return an error if the role ID is not found in the system roles
}
