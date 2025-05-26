package http

import (
	"delivery/internal/adapters/in/http/problems"
	"delivery/internal/core/application/usecases/queries"
	"delivery/internal/generated/servers"
	"delivery/internal/pkg/errs"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *Server) GetCouriers(ctx echo.Context) error {
	query, err := queries.NewGetAllCouriersQuery()
	if err != nil {
		return problems.NewBadRequest(err.Error())
	}

	response, err := s.GetAllCouriersQueryHandler().Handle(query)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return ctx.JSON(http.StatusNotFound, problems.NewNotFound(err.Error()))
		}
	}

	couriers := mapToCouriersDto(response)
	return ctx.JSON(http.StatusOK, couriers)
}

func mapToCouriersDto(response queries.GetAllCouriersResponse) []servers.Courier {
	var couriers []servers.Courier
	for _, courier := range response.Couriers {
		location := servers.Location{
			X: courier.Location.X,
			Y: courier.Location.Y,
		}

		var courier = servers.Courier{
			Id:       courier.ID,
			Name:     courier.Name,
			Location: location,
		}
		couriers = append(couriers, courier)
	}
	return couriers
}
