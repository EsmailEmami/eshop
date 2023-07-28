package errors

import (
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ValidationErrors struct {
	Errs map[string][]string `json:"errors,omitempty"`
}

func (verr ValidationErrors) Error() string {
	errString := ""
	for k, errs := range verr.Errs {
		errString += k + ":" + strings.Join(errs, ";")
	}
	return errString
}

type ValidationError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
	ValidationErrors
}

func (e *ValidationError) Error() string {
	return e.Message
}

func NewValidationError(message string, err error) *ValidationError {
	if err == nil {
		return &ValidationError{Message: message, StatusCode: http.StatusUnprocessableEntity}
	}
	switch e := err.(type) {
	case validation.Errors:
		{
			verr := &ValidationError{Message: message, StatusCode: http.StatusUnprocessableEntity}
			verr.ValidationErrors = ValidationErrors{}
			verr.ValidationErrors.Errs = make(map[string][]string)
			for k, v := range e {
				verr.ValidationErrors.Errs[k] = append(verr.ValidationErrors.Errs[k], v.Error())
			}
			return verr
		}
	case ValidationErrors:
		{
			verr := &ValidationError{
				Message:          message,
				StatusCode:       http.StatusUnprocessableEntity,
				ValidationErrors: e,
			}

			return verr
		}
	default:
		{
			return &ValidationError{Message: message, StatusCode: http.StatusUnprocessableEntity}
		}
	}
}

type BadRequestError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
	LogError   error  `json:"-"`
}

func (e *BadRequestError) Error() string {
	return e.Message
}

func NewBadRequestError(message string, logError error) *BadRequestError {
	return &BadRequestError{Message: message, StatusCode: http.StatusBadRequest, LogError: logError}
}

type RecordNotFoundError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
	LogError   error  `json:"-"`
}

func (e *RecordNotFoundError) Error() string {
	return e.Message
}

func NewRecordNotFoundError(message string, logError error) *RecordNotFoundError {
	return &RecordNotFoundError{Message: message, StatusCode: http.StatusNotFound, LogError: logError}
}

type InternalServerError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
	LogError   error  `json:"-"`
}

func (e *InternalServerError) Error() string {
	return e.Message
}

func NewInternalServerError(message string, logError error) *InternalServerError {
	return &InternalServerError{Message: message, StatusCode: http.StatusInternalServerError, LogError: logError}
}

type UnauthorizedError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
	LogError   error  `json:"-"`
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

func NewUnauthorizedError(message string, logError error) *UnauthorizedError {
	return &UnauthorizedError{Message: message, StatusCode: http.StatusUnauthorized, LogError: logError}
}

type ForbiddenError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
	LogError   error  `json:"-"`
}

func (e *ForbiddenError) Error() string {
	return e.Message
}

func NewForbiddenError(message string, logError error) *ForbiddenError {
	return &ForbiddenError{Message: message, StatusCode: http.StatusForbidden, LogError: logError}
}
