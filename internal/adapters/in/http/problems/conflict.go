package problems

import (
	"errors"
	"net/http"
)

var ProblemConflict = errors.New("conflict")

type ConflictError struct {
	ProblemDetails
}

func NewConflict(problemType string, detail string) *ConflictError {
	return &ConflictError{
		ProblemDetails: ProblemDetails{
			Type:   problemType,
			Title:  "Conflict",
			Status: http.StatusConflict,
			Detail: detail,
		},
	}
}

func (e *ConflictError) Error() string {
	return e.ProblemDetails.Error()
}

func (e *ConflictError) Unwrap() error {
	return ProblemConflict
}
