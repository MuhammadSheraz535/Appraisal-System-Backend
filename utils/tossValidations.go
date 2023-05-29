package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/mrehanabbasi/appraisal-system-backend/constants"
	log "github.com/mrehanabbasi/appraisal-system-backend/logger"
)

func CheckKpiAgainstTossApis(selectedAssignID uint16, assignType string) (int, error) {
	// Check which SelectedAssignID exists in the API
	tossBaseUrl := os.Getenv("TOSS_BASE_URL")

	switch assignType {
	case constants.ASSIGN_TYPE_ROLE:
		method := http.MethodGet
		url := tossBaseUrl + "/api/Employee/GetSystemRolesList"

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

		type SystemRole struct {
			Value uint16 `json:"value"`
			Label string `json:"label"`
		}

		type SystemRolesResponse struct {
			SystemRoles []SystemRole `json:"systemRoles"`
		}

		var response SystemRolesResponse
		if err := json.Unmarshal(responseBody, &response); err != nil {
			log.Error(err.Error())
			return http.StatusInternalServerError, err
		}

		found := false
		for _, role := range response.SystemRoles {
			if role.Value == selectedAssignID {
				found = true
				break
			}
		}

		if !found {
			err := errors.New("invalid selected role id")
			log.Error(err.Error())
			return http.StatusBadRequest, err
		}
	case constants.ASSIGN_TYPE_TEAM:
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
				ProjectID uint16 `json:"projectId"`
			} `json:"projectDetails"`
		}

		if err := json.Unmarshal(responseBody, &projects); err != nil {
			log.Error(err.Error())
			return http.StatusInternalServerError, err
		}

		found := false
		for _, project := range projects {
			if project.ProjectDetails.ProjectID == selectedAssignID {
				found = true
				break
			}
		}

		if !found {
			err := errors.New("invalid selected team id")
			log.Error(err.Error())
			return http.StatusBadRequest, err
		}
	case constants.ASSIGN_TYPE_INDIVIDUAL:
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
			if employee.EmployeeID == selectedAssignID {
				found = true
				break
			}
		}

		if !found {
			err := errors.New("invalid selected employee id")
			log.Error(err.Error())
			return http.StatusBadRequest, err
		}
	}

	return 0, nil
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

func CheckTeamAgainstToss(TeamId uint16) (int, error) {

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
			ProjectID uint16 `json:"projectId"`
		} `json:"projectDetails"`
	}

	if err := json.Unmarshal(responseBody, &projects); err != nil {
		log.Error(err.Error())
		return http.StatusInternalServerError, err
	}

	found := false
	for _, project := range projects {
		if project.ProjectDetails.ProjectID == TeamId {
			found = true
			break
		}
	}

	if !found {
		err := errors.New("invalid selected team id")
		log.Error(err.Error())
		return http.StatusBadRequest, err
	}
	return 0, nil
}
