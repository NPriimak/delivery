package cmd

import (
	"delivery/internal/adapters/out/postgres"
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/core/application/usecases/queries"
	"delivery/internal/core/domain/services"
	"delivery/internal/core/ports"
	"gorm.io/gorm"
)

type CompositionRoot struct {
	configs Config
	gormDb  *gorm.DB

	closers []Closer
}

func NewCompositionRoot(c Config, gormDb *gorm.DB) CompositionRoot {
	app := CompositionRoot{
		configs: c,
	}
	return app
}

func (cr *CompositionRoot) NewAssignOrdersCommandHandler() commands.AssignOrderCommandHandler {
	handler, err := commands.NewAssignOrderCommandHandler(cr.NewUnitOfWork(), cr.NewOrderDispatcher())
	if err != nil {
		panic(err)
	}
	return handler
}

func (cr *CompositionRoot) NewCreateCourierCommandHandler() commands.CreateCourierCommandHandler {
	handler, err := commands.NewCreateCourierCommandHandler(cr.NewUnitOfWork())
	if err != nil {
		panic(err)
	}
	return handler
}

func (cr *CompositionRoot) NewCreateOrderCommandHandler() commands.CreateOrderCommandHandler {
	handler, err := commands.NewCreateOrderCommandHandler(cr.NewUnitOfWork())
	if err != nil {
		panic(err)
	}
	return handler
}

func (cr *CompositionRoot) NewMoveCouriersCommandHandler() commands.MoveCouriersCommandHandler {
	handler, err := commands.NewMoveCouriersCommandHandler(cr.NewUnitOfWork())
	if err != nil {
		panic(err)
	}
	return handler
}

func (cr *CompositionRoot) NewGetAllCouriersQueryHandler() queries.GetAllCouriersQueryHandler {
	handler, err := queries.NewGetAllCouriersQueryHandler(cr.gormDb)
	if err != nil {
		panic(err)
	}
	return handler
}

func (cr *CompositionRoot) NewGetNotCompletedOrdersQueryHandler() queries.GetNotCompletedOrdersQueryHandler {
	handler, err := queries.NewGetNotCompletedOrdersQueryHandler(cr.gormDb)
	if err != nil {
		panic(err)
	}
	return handler
}

func (cr *CompositionRoot) NewUnitOfWork() ports.UnitOfWork {
	uow, err := postgres.NewUnitOfWork(cr.gormDb)
	if err != nil {
		panic(err)
	}
	return uow
}

func (cr *CompositionRoot) NewOrderDispatcher() services.OrderDispatcher {
	return services.NewOrderDispatcher()
}
