package models

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

type TokenInfo struct {
	EmpID         uint16
	EmailAddr     string
	DesignationID uint16
	Designation   string
	SupervisorID  uint16
	Supervisor    string
	Department    string
	EmpImagePath  string
	EmpRoleID     uint16
}

type TossClaims struct {
	jwt.RegisteredClaims
	ClientID         string           `json:"client_id"`
	AuthTime         *jwt.NumericDate `json:"auth_time"`
	IdentityProvider string           `json:"idp"`
	EmpID            string           `json:"id"`
	Locale           string           `json:"locale"`
	Email            string           `json:"email"`
	EmailAddr        string           `json:"emailAdd"`
	Designation      string           `json:"designation"`
	EmpDesignation   string           `json:"empDesignation"`
	Supervisor       string           `json:"supervisor"`
	Department       string           `json:"department"`
	JoiningDate      string           `json:"empJoinedOn"`
	EmpImagePath     string           `json:"empImagePath"`
	Role             string           `json:"role"`
	EmpRole          string           `json:"empRole"`
	Scope            jwt.ClaimStrings `json:"scope"`
	Amr              jwt.ClaimStrings `json:"amr"`
}

func (t TossClaims) Validate(ctx context.Context) error {
	if t.ExpiresAt == nil {
		return errors.New("exp not provided in the token")
	}
	if t.ClientID == "" {
		return errors.New("client_id not provided in the token")
	}
	if t.EmpID == "" {
		return errors.New("id not provided in the token")
	}
	if t.Email == "" {
		return errors.New("email not provided in the token")
	}
	if t.EmailAddr == "" {
		return errors.New("emailAdd not provided in the token")
	}
	if t.Designation == "" {
		return errors.New("designation not provided in the token")
	}
	if t.EmpDesignation == "" {
		return errors.New("empDesignation not provided in the token")
	}
	if t.EmpDesignation == "" {
		return errors.New("empDesignation not provided in the token")
	}
	if t.Supervisor == "" {
		return errors.New("supervisor not provided in the token")
	}
	if t.Department == "" {
		return errors.New("department not provided in the token")
	}
	if t.JoiningDate == "" {
		return errors.New("empJoinedOn not provided in the token")
	}
	if t.EmpImagePath == "" {
		return errors.New("empImagePath not provided in the token")
	}
	if t.Role == "" {
		return errors.New("role not provided in the token")
	}
	if t.EmpRole == "" {
		return errors.New("empRole not provided in the token")
	}
	return nil
}
