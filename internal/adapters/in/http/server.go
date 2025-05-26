package http

import (
	"delivery/cmd"
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/core/application/usecases/queries"
	"delivery/internal/generated/servers"
)

var _ servers.ServerInterface = &Server{}

type Server struct {
	root cmd.CompositionRoot
}

func NewServer(root cmd.CompositionRoot) (*Server, error) {
	return &Server{root}, nil
}

func (s *Server) CreateOrderCommandHandler() commands.CreateOrderCommandHandler {
	return s.root.NewCreateOrderCommandHandler()
}

func (s *Server) CreateCourierCommandHandler() commands.CreateCourierCommandHandler {
	return s.root.NewCreateCourierCommandHandler()
}

func (s *Server) GetAllCouriersQueryHandler() queries.GetAllCouriersQueryHandler {
	return s.root.NewGetAllCouriersQueryHandler()
}

func (s *Server) GetNotCompletedOrdersQueryHandler() queries.GetNotCompletedOrdersQueryHandler {
	return s.root.NewGetNotCompletedOrdersQueryHandler()
}
