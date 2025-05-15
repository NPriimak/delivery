package problems

import (
	"errors"
	"net/http"
)

var ProblemBadRequest = errors.New("bad request")

type BadRequest struct {
	ProblemDetails
}

func NewBadRequest(detail string) *BadRequest {
	return &BadRequest{
		ProblemDetails: ProblemDetails{
			Type:   "bad-request",
			Title:  "Bad Request",
			Status: http.StatusBadRequest,
			Detail: detail,
		},
	}
}

func (e *BadRequest) Error() string {
	return e.ProblemDetails.Error()
}

func (e *BadRequest) Unwrap() error {
	return ProblemBadRequest
}
