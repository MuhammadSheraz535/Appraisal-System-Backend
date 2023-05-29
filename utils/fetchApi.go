package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

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
			RoleID uint16 `json:"empRole"` // Struct for unmarshaling the JSON response
		}

		if err := json.Unmarshal(responseBody, &employee); err != nil {
			log.Error(err.Error())
			return nil, err // Return an error if there's an issue unmarshaling the JSON response
		}

		roleIDs = append(roleIDs, employee.RoleID) // Append the extracted role ID to the slice
	}

	return roleIDs, nil // Return the list of role IDs
}

func GetEmployeesId(teamID uint16) ([]uint16, error) {
	type ProjectResponse struct {
		ProjectDetails struct {
			ProjectID    uint16 `json:"projectId"`
			SupervisorID uint16 `json:"supervisorId"`
		} `json:"projectDetails"`
		AllocateTo []struct {
			EmployeeID uint16 `json:"employeeId"`
		} `json:"allocateTo"`
	}

	tossBaseUrl := os.Getenv("TOSS_BASE_URL")          // Get the TOSS base URL from environment variable
	method := http.MethodGet                           // HTTP method for sending the request
	url := tossBaseUrl + "/api/Project/GetAllProjects" // Construct the URL for fetching all projects

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

	var projects []ProjectResponse // Slice to store the unmarshaled project responses

	if err := json.Unmarshal(responseBody, &projects); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	var employeeIDs []uint16

	for _, project := range projects {
		if project.ProjectDetails.ProjectID == teamID { // Check if the project ID matches the provided team ID
			for _, allocateTo := range project.AllocateTo {
				if allocateTo.EmployeeID != project.ProjectDetails.SupervisorID {
					employeeIDs = append(employeeIDs, allocateTo.EmployeeID) // Append employee ID to the slice if it's not the supervisor
				}
			}
			break // Break the loop as we found the matching project
		}
	}

	return employeeIDs, nil // Return the list of employee IDs
}

func VerifyTeamAndSupervisorID(teamID, supervisorID uint16) (int, error) {
	tossBaseUrl := os.Getenv("TOSS_BASE_URL")

	method := http.MethodGet
	url := tossBaseUrl + "/api/Project/GetAllProjects"

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

	var projects []struct {
		ProjectDetails struct {
			ProjectID    uint16 `json:"projectId"`
			SupervisorID uint16 `json:"supervisorId"`
		} `json:"projectDetails"`
	}

	if err := json.Unmarshal(responseBody, &projects); err != nil {
		log.Error(err.Error())
		return http.StatusInternalServerError, err
	}

	foundteam, foundsupervisor := false, false
	for _, project := range projects {
		if project.ProjectDetails.ProjectID == teamID {
			foundteam = true
			if project.ProjectDetails.SupervisorID == supervisorID {

				foundsupervisor = true
				break
			}
			break
		}
	}

	if !foundteam {
		err := errors.New("invalid selected team id")
		log.Error(err.Error())
		return http.StatusBadRequest, err
	}
	if !foundsupervisor {
		err := errors.New("invalid selected supervisor id")
		log.Error(err.Error())
		return http.StatusBadRequest, err
	}
	return 0, nil
}
