package services

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type (
    ErrorResponse struct {
        Error       bool
        FailedField string
        Tag         string
        Value       interface{}
    }

    XValidator struct {
        validator *validator.Validate
    }

    GlobalErrorHandlerResp struct {
        Success bool   `json:"success"`
        Message string `json:"message"`
    }
)

var validate = validator.New()

func (v XValidator) Validate(data interface{}) []ErrorResponse {
    validationErrors := []ErrorResponse{}

    errs := validate.Struct(data)
    if errs != nil {
        for _, err := range errs.(validator.ValidationErrors) {
            var elem ErrorResponse

            elem.FailedField = err.Field()
            elem.Tag = err.Tag()
            elem.Value = err.Value()
            elem.Error = true

            validationErrors = append(validationErrors, elem)
        }
    }

    return validationErrors
}

func FormatErrors(errs []ErrorResponse) []string {
    errMsgs := make([]string, 0)

    for _, err := range errs {
        errStr := "";

        switch err.Tag {
            case "required":
                errStr = fmt.Sprintf("%s is required", err.FailedField)
                break;
            case "min":
                errStr = fmt.Sprintf("%s is too short", err.FailedField)
                break;
            case "max":
                errStr = fmt.Sprintf("%s is too long", err.FailedField)
                break;
            default:
                errStr = fmt.Sprintf("%s is '%v' but needs to implement '%s'", err.FailedField, err.Value, err.Tag)
        }

        errMsgs = append(errMsgs, errStr)
    }

    return errMsgs
}

func GetValidator() *XValidator {
    return &XValidator{
        validator: validate,
    }
}