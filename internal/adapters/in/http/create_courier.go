package http

import (
	"delivery/internal/adapters/in/http/problems"
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/generated/servers"
	"delivery/internal/pkg/errs"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *Server) CreateCourier(ctx echo.Context) error {
	var courier servers.NewCourier
	if err := ctx.Bind(&courier); err != nil {
		return problems.NewBadRequest("invalid JSON body: " + err.Error())
	}

	createCourierCommand, err := commands.NewCreateCourierCmd(courier.Name, courier.Speed)
	if err != nil {
		return problems.NewBadRequest(err.Error())
	}

	err = s.CreateCourierCommandHandler().Handle(ctx.Request().Context(), createCourierCommand)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return problems.NewNotFound(err.Error())
		}
		return problems.NewConflict(err.Error(), "/")
	}

	return ctx.JSON(http.StatusOK, nil)
}
