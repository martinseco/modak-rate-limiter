package errors

import (
	"fmt"
	"net/http"
)

type (
	BadRequestError string
	ForbiddenError  string
	NotFoundError   string

	ApiError struct {
		ErrorStr string   `json:"error,omitempty"`
		Message  string   `json:"message,omitempty"`
		Status   int      `json:"status,omitempty"`
		Causes   []string `json:"causes,omitempty"`
	}
)

func (e BadRequestError) Error() string { return string(e) }
func (e ForbiddenError) Error() string  { return string(e) }
func (e NotFoundError) Error() string   { return string(e) }
func (e ApiError) Error() string        { return fmt.Sprintf("%s - %s", e.ErrorStr, e.Message) }

func badRequestApiError(message string) ApiError {
	return ApiError{ErrorStr: http.StatusText(http.StatusBadRequest), Message: message, Status: http.StatusBadRequest}
}

func forbiddenApiError(message string) ApiError {
	return ApiError{ErrorStr: http.StatusText(http.StatusForbidden), Message: message, Status: http.StatusForbidden}
}

func notFoundApiError(message string) ApiError {
	return ApiError{ErrorStr: http.StatusText(http.StatusNotFound), Message: message, Status: http.StatusNotFound}
}

func internalApiError(err error) ApiError {
	return ApiError{ErrorStr: http.StatusText(http.StatusInternalServerError), Message: err.Error(), Status: http.StatusInternalServerError}
}
