package util

import (
	"github.com/go-playground/validator/v10"
	"regexp"
	"strings"
)

type Validator interface {
	Validate(i interface{}) error
	GetErrorMessage(err validator.ValidationErrors) string
}

type customValidator struct {
	Validator *validator.Validate
}

func GetNewValidator() *customValidator {
	validate := validator.New()

	err := validate.RegisterValidation("alphanumeric", validateAlphaNumericVal)
	if err != nil {
		return nil
	}

	return &customValidator{Validator: validate}
}

func validateAlphaNumericVal(fl validator.FieldLevel) bool {
	re := regexp.MustCompile("^[ a-zA-Z0-9_-]*$")
	return re.MatchString(fl.Field().String())
}

func (cv *customValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

func (cv *customValidator) GetErrorMessage(err validator.ValidationErrors) string {
	if err == nil || len(err) == 0 {
		return ""
	}

	var sb strings.Builder

	sb.WriteString("validation failed on field '" + err[0].Field() + "'")

	return sb.String()
}
