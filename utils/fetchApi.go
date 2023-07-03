package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	log "github.com/mrehanabbasi/appraisal-system-backend/logger"
)

func GetRolesID(empIds []uint16) ([]uint16, error) {
	tossBaseUrl := os.Getenv("TOSS_BASE_URL") // Get the TOSS base URL from environment variable
	method := http.MethodGet                  // HTTP method for sending the request

	var roleIDs []uint16

	for _, empId := range empIds {
		url := tossBaseUrl + "/api/Employee/" + strconv.FormatUint(uint64(empId), 10) // Construct the URL for fetching employee details based on empId

		resp, err := SendRequest(method, url, nil) // Send the HTTP request to the specified URL
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			errMsg := fmt.Sprintf("failed to get role id against employee id %d. status code: %d", empId, resp.StatusCode)
			log.Error(errMsg)
			return nil, errors.New(errMsg) // Return an error if the response status code is not OK
		}

		responseBody, err := io.ReadAll(resp.Body) // Read the response body
		if err != nil {
			log.Error(err.Error())
			return nil, err // Return an error if there's an issue reading the response body
		}

		var employee struct {
			RoleID uint16 `json:"empDesignation"` // Struct for unmarshaling the JSON response
		}

		if err := json.Unmarshal(responseBody, &employee); err != nil {
			log.Error(err.Error())
			return nil, err // Return an error if there's an issue unmarshaling the JSON response
		}

		roleIDs = append(roleIDs, employee.RoleID) // Append the extracted role ID to the slice
	}

	return roleIDs, nil // Return the list of role IDs
}

type Employee struct {
	ID                        uint16 `json:"id"`
	EmployeeID                uint16 `json:"employeeId"`
	EmployeeName              string `json:"employeeName"`
	ProjectName               string `json:"projectName"`
	ProjectStartedDate        string `json:"projectStartedDate"`
	EmployeeProjectSupervisor string `json:"employeeProjectSupervisor"`
}

type ProjectResponse struct {
	ProjectID        uint16     `json:"projectId"`
	ProjectName      string     `json:"projectName"`
	ProjectEmployees []Employee `json:"projectEmployees"`
}

func GetEmployeesId(teamID uint16) ([]uint16, error) {

	tossBaseUrl := os.Getenv("TOSS_BASE_URL")
	method := http.MethodGet
	url := tossBaseUrl + "/api/Project/AllProjectsWithEmployeesList?IsActive=true"

	resp, err := SendRequest(method, url, nil)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	var projects []ProjectResponse

	if err := json.Unmarshal(responseBody, &projects); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	var employeeIDs []uint16

	for _, project := range projects {
		if project.ProjectID == teamID {
			for _, employee := range project.ProjectEmployees {
				if employee.EmployeeProjectSupervisor != employee.EmployeeName {
					employeeIDs = append(employeeIDs, employee.EmployeeID)
				}
			}
			break
		}
	}

	return employeeIDs, nil
}

func VerifyTeamAndSupervisorID(teamID, supervisorID uint16) (int, string, error) {
	tossBaseUrl := os.Getenv("TOSS_BASE_URL")

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

	var projects []ProjectResponse

	if err := json.Unmarshal(responseBody, &projects); err != nil {
		log.Error(err.Error())
		return http.StatusInternalServerError, "", err
	}

	var teamName string
	foundTeam, foundSupervisor := false, false

	for _, project := range projects {
		if project.ProjectID == teamID {
			foundTeam = true
			for _, employee := range project.ProjectEmployees {
				if employee.EmployeeProjectSupervisor == employee.EmployeeName && employee.EmployeeID == supervisorID {
					foundSupervisor = true
					teamName = project.ProjectName
					break
				}
			}
			break
		}
	}

	if !foundTeam {
		err := errors.New("invalid selected team ID")
		log.Error(err.Error())
		return http.StatusBadRequest, "", err
	}

	if !foundSupervisor {
		err := errors.New("invalid selected supervisor ID")
		log.Error(err.Error())
		return http.StatusBadRequest, "", err
	}
	trimmedStr := strings.Trim(teamName, "\r\n")
	return 0, trimmedStr, nil
}

func GetSupervisorName(SprID uint16) (string, error) {
	tossBaseUrl := os.Getenv("TOSS_BASE_URL") // Get the TOSS base URL from the environment variable
	method := http.MethodGet                  // HTTP method for sending the request

	url := tossBaseUrl + "/api/Project/AllProjectsWithEmployeesList?IsActive=true" // Construct the URL for fetching all projects

	resp, err := SendRequest(method, url, nil) // Send the HTTP request to the specified URL
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg := "Failed to get supervisor name for supervisor ID: " + strconv.Itoa(int(SprID)) + ". Status code: " + strconv.Itoa(resp.StatusCode)
		return "", errors.New(errMsg) // Return an error if the response status code is not OK
	}

	var projects []ProjectResponse

	responseBody, err := io.ReadAll(resp.Body) // Read the response body
	if err != nil {
		return "", err // Return an error if there's an issue reading the response body
	}

	if err := json.Unmarshal(responseBody, &projects); err != nil {
		return "", err // Return an error if there's an issue unmarshaling the JSON response
	}

	for _, project := range projects {
		for _, employee := range project.ProjectEmployees {
			if employee.EmployeeID == SprID && employee.EmployeeName == employee.EmployeeProjectSupervisor {
				return employee.EmployeeName, nil // Return the Supervisor Name if the supervisor ID and name match
			}
		}
	}

	return "", errors.New("supervisor not found") // Return an error if the supervisor ID is not found in the projects
}

func VerifyIndividualAndSupervisorID(indID, supervisorID uint16) (int, string, error) {
	tossBaseUrl := os.Getenv("TOSS_BASE_URL")
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

	var projects []ProjectResponse

	if err := json.Unmarshal(responseBody, &projects); err != nil {
		log.Error(err.Error())
		return http.StatusInternalServerError, "", err
	}

	var empName string
	foundIndividual, foundSupervisor := false, false

	for _, project := range projects {
		for _, employee := range project.ProjectEmployees {
			if employee.EmployeeID == supervisorID && employee.EmployeeName == employee.EmployeeProjectSupervisor {
				empName = employee.EmployeeProjectSupervisor // Return the Supervisor Name if the supervisor ID and name match
				foundSupervisor = true
			}
			if employee.EmployeeID == indID && empName == employee.EmployeeProjectSupervisor {
				empName = employee.EmployeeName // Return the Individual's Name if the individual ID and supervisor name match
				foundIndividual = true
				break
			}

		}
	}

	if !foundSupervisor {
		err := errors.New("invalid selected supervisor ID")
		log.Error(err.Error())
		return http.StatusBadRequest, "", err
	}

	if !foundIndividual {
		err := errors.New("invalid selected individual ID")
		log.Error(err.Error())
		return http.StatusBadRequest, "", err
	}

	return 0, empName, nil
}

func GetDesignationName(DesignationID uint16) (string, error) {
	tossBaseUrl := os.Getenv("TOSS_BASE_URL") // Get the TOSS base URL from environment variable
	method := http.MethodGet                  // HTTP method for sending the request

	url := tossBaseUrl + "/api/Designation/" + strconv.FormatUint(uint64(DesignationID), 10) // Construct the URL for fetching designation details based on designation ID

	resp, err := SendRequest(method, url, nil) // Send the HTTP request to the specified URL
	if err != nil {
		log.Error(err.Error())
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("failed to get designation name for designation ID %d. status code: %d", DesignationID, resp.StatusCode)
		log.Error(errMsg)
		return "", errors.New(errMsg) // Return an error if the response status code is not OK
	}

	responseBody, err := io.ReadAll(resp.Body) // Read the response body
	if err != nil {
		log.Error(err.Error())
		return "", err // Return an error if there's an issue reading the response body
	}

	var designation struct {
		DesignationID   uint16 `json:"designationId"`
		DesignationName string `json:"designationName"` // Struct for unmarshaling the JSON response
	}

	if err := json.Unmarshal(responseBody, &designation); err != nil {
		log.Error(err.Error())
		return "", err // Return an error if there's an issue unmarshaling the JSON response
	}

	return designation.DesignationName, nil // Return the designation name
}

func GetEmployeeIDsByDesignation(designation uint16) ([]uint16, error) {
	tossBaseURL := os.Getenv("TOSS_BASE_URL") // Get the TOSS base URL from the environment variable
	method := http.MethodGet                  // HTTP method for sending the request

	url := tossBaseURL + "/api/Employee/GetAllEmployeesInfo?Status=1&PageSize=10000&Designation=" + strconv.Itoa(int(designation))

	resp, err := SendRequest(method, url, nil) // Send the HTTP request to the specified URL
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("failed to get employee IDs for designation '%d'. status code: %d", designation, resp.StatusCode)
		log.Error(errMsg)
		return nil, errors.New(errMsg) // Return an error if the response status code is not OK
	}

	responseBody, err := io.ReadAll(resp.Body) // Read the response body
	if err != nil {
		log.Error(err.Error())
		return nil, err // Return an error if there's an issue reading the response body
	}

	type EmployeeInfo struct {
		ID   uint16 `json:"id"`
		Name string `json:"name"`
	}

	var response struct {
		EmployeeInfo []EmployeeInfo `json:"employeeInfo"`
	}

	if err := json.Unmarshal(responseBody, &response); err != nil {
		log.Error(err.Error())
		return nil, err // Return an error if there's an issue unmarshaling the JSON response
	}

	employeeIDs := make([]uint16, 0, len(response.EmployeeInfo))
	for _, employee := range response.EmployeeInfo {
		employeeIDs = append(employeeIDs, employee.ID) // Append the employee ID to the slice
	}

	return employeeIDs, nil // Return the list of employee IDs
}
func GetProjectDetailsByEmployeeID(employeeID uint16) ([]ProjectResponse, error) {
	tossBaseUrl := os.Getenv("TOSS_BASE_URL") // Get the TOSS base URL from environment variable
	method := http.MethodGet                  // HTTP method for sending the request
	url := tossBaseUrl + "/api/Project/AllProjectsWithEmployeesList?IsActive=true"

	resp, err := SendRequest(method, url, nil) // Send the HTTP request to fetch all projects
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	var projects []ProjectResponse

	if err := json.Unmarshal(responseBody, &projects); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	var projectDetails []ProjectResponse

	for _, project := range projects {
		for _, employee := range project.ProjectEmployees {
			if employee.EmployeeID == employeeID {
				projectDetails = append(projectDetails, project)
				break
			}
		}
	}

	return projectDetails, nil
}

type EmployeeImage struct {
	EmployeeImage string `json:"employeeImage"`
}

func GetEmployeeImageByID(employeeID uint64) (string, error) {

	tossBaseURL := os.Getenv("TOSS_BASE_URL") // Get the TOSS base URL from the environment variable
	method := http.MethodGet                  // HTTP method for sending the request

	url := tossBaseURL + "/api/Employee/" + strconv.FormatUint(uint64(employeeID), 10)

	resp, err := SendRequest(method, url, nil) // Send the HTTP request to fetch all projects
	if err != nil {
		log.Error(err.Error())
		return "", err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err.Error())
		return "", err
	}

	var employeeImage EmployeeImage
	if err := json.Unmarshal(responseBody, &employeeImage); err != nil {
		return "", err
	}

	return employeeImage.EmployeeImage, nil
}
