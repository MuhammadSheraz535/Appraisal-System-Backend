package controller

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ErrValidationSlice(err error) ([]string, bool) {
	errs := []string{}
	ve, ok := err.(validator.ValidationErrors)
	if ok {
		for _, fe := range ve {
			errStr := errMsgFromValidator(fe.Field(), fe.Tag(), fe.Param())
			errs = append(errs, errStr)
		}
	}

	return errs, ok
}

func errMsgFromValidator(field, tag, value string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("%s field is required", field)
	case "email":
		return "invalid email"
	case "min":
		return fmt.Sprintf("minimum %v characters required for %s", value, field)
	case "max":
		return fmt.Sprintf("maximum %v characters allowed for %s", value, field)
	}
	return tag
}
