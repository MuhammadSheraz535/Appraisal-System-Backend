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
		url := tossBaseUrl + "/api/Project/AllProjectsWithEmployeesList?IsActive=true"

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
		type ProjectResponse struct {
			ProjectID   uint16 `json:"projectId"`
			ProjectName string `json:"projectName"`
		}

		var projects []ProjectResponse

		if err := json.Unmarshal(responseBody, &projects); err != nil {
			log.Error(err.Error())
			return http.StatusInternalServerError, "", err
		}

		found := false
		for _, project := range projects {
			if project.ProjectID == selectedAssignID {
				found = true
				selectedAssignName = project.ProjectName
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
