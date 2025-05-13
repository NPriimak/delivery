package http

import (
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/core/application/usecases/queries"
	"delivery/internal/generated/servers"
	"delivery/internal/pkg/errs"
)

var _ servers.ServerInterface = &Server{}

type Server struct {
	createOrderCommandHandler   commands.CreateOrderCommandHandler
	createCourierCommandHandler commands.CreateCourierCommandHandler

	getAllCouriersQueryHandler        queries.GetAllCouriersQueryHandler
	getNotCompletedOrdersQueryHandler queries.GetNotCompletedOrdersQueryHandler
}

func NewServer(
	createOrderCommandHandler commands.CreateOrderCommandHandler,
	createCourierCommandHandler commands.CreateCourierCommandHandler,

	getAllCouriersQueryHandler queries.GetAllCouriersQueryHandler,
	getNotCompletedOrdersQueryHandler queries.GetNotCompletedOrdersQueryHandler,
) (*Server, error) {
	if createOrderCommandHandler == nil {
		return nil, errs.NewValueIsRequiredError("createOrderCommandHandler")
	}
	if createCourierCommandHandler == nil {
		return nil, errs.NewValueIsRequiredError("createCourierCommandHandler")
	}
	if getAllCouriersQueryHandler == nil {
		return nil, errs.NewValueIsRequiredError("getAllCouriersQueryHandler")
	}
	if getNotCompletedOrdersQueryHandler == nil {
		return nil, errs.NewValueIsRequiredError("getNotCompletedOrdersQueryHandler")
	}
	return &Server{
		createOrderCommandHandler:         createOrderCommandHandler,
		createCourierCommandHandler:       createCourierCommandHandler,
		getAllCouriersQueryHandler:        getAllCouriersQueryHandler,
		getNotCompletedOrdersQueryHandler: getNotCompletedOrdersQueryHandler,
	}, nil
}
