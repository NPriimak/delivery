package problems

import (
	"errors"
	"net/http"
)

var NotFound = errors.New("not found")

type NotFoundError struct {
	ProblemDetails
}

func NewNotFound(detail string) *NotFoundError {
	return &NotFoundError{
		ProblemDetails: ProblemDetails{
			Type:   "not-found",
			Title:  "Resource Not Found",
			Status: http.StatusNotFound,
			Detail: detail,
		},
	}
}

func (e *NotFoundError) Error() string {
	return e.ProblemDetails.Error()
}

func (e *NotFoundError) Unwrap() error {
	return NotFound
}
