package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	log "github.com/mrehanabbasi/appraisal-system-backend/logger"
)

func GetRolesID(empIds []uint16) ([]uint16, error) {
	tossBaseUrl := os.Getenv("TOSS_BASE_URL")
	method := http.MethodGet

	var roleIDs []uint16

	for _, empId := range empIds {
		url := tossBaseUrl + "/api/Employee/" + strconv.FormatUint(uint64(empId), 10)

		resp, err := SendRequest(method, url, nil)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			errMsg := fmt.Sprintf("failed to get role id against employee id %d. status code: %d", empId, resp.StatusCode)
			log.Error(errMsg)
			return nil, errors.New(errMsg)
		}

		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}

		var employee struct {
			RoleID uint16 `json:"empRole"`
		}

		if err := json.Unmarshal(responseBody, &employee); err != nil {
			log.Error(err.Error())
			return nil, err
		}

		roleIDs = append(roleIDs, employee.RoleID)
	}

	return roleIDs, nil
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

	tossBaseUrl := os.Getenv("TOSS_BASE_URL")
	method := http.MethodGet
	url := tossBaseUrl + "/api/Project/GetAllProjects"

	resp, err := SendRequest(method, url, nil)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
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
		if project.ProjectDetails.ProjectID == teamID {
			for _, allocateTo := range project.AllocateTo {
				if allocateTo.EmployeeID != project.ProjectDetails.SupervisorID {
					employeeIDs = append(employeeIDs, allocateTo.EmployeeID)
				}
			}
			break
		}
	}

	return employeeIDs, nil
}
