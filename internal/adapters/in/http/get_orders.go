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

func (s *Server) GetOrders(c echo.Context) error {
	query, err := queries.NewGetNotCompletedOrdersQuery()
	if err != nil {
		return problems.NewBadRequest(err.Error())
	}

	response, err := s.getNotCompletedOrdersQueryHandler.Handle(query)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return c.JSON(http.StatusNotFound, problems.NewNotFound(err.Error()))
		}
	}

	orders := mapToOrdersDto(response)
	return c.JSON(http.StatusOK, orders)
}

func mapToOrdersDto(response queries.GetNotCompletedOrdersResponse) []servers.Order {
	var orders []servers.Order
	for _, courier := range response.Orders {
		location := servers.Location{
			X: courier.Location.X,
			Y: courier.Location.Y,
		}

		var courier = servers.Order{
			Id:       courier.ID,
			Location: location,
		}
		orders = append(orders, courier)
	}
	return orders
}
