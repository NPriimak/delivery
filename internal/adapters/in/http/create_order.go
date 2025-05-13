package http

import (
	"delivery/internal/adapters/in/http/problems"
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/pkg/errs"
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *Server) CreateOrder(c echo.Context) error {
	createOrderCommand, err := commands.NewCreateOrderCmd(uuid.New(), "Street", 5)
	if err != nil {
		return problems.NewBadRequest(err.Error())
	}

	err = s.createOrderCommandHandler.Handle(c.Request().Context(), createOrderCommand)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return problems.NewNotFound(err.Error())
		}
		return problems.NewConflict(err.Error(), "/")
	}

	return c.JSON(http.StatusOK, nil)
}
