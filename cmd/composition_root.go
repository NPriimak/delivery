package cmd

import (
	"delivery/internal/adapters/in/jobs"
	"delivery/internal/adapters/out/postgres/courierrepo"
	"delivery/internal/adapters/out/postgres/orderrepo"
	"delivery/internal/adapters/out/postgres/shared"
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/core/application/usecases/queries"
	"delivery/internal/core/domain/services"
	"delivery/internal/core/ports"
	"github.com/robfig/cron/v3"
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
		gormDb:  gormDb,
	}
	return app
}

func (cr *CompositionRoot) NewAssignOrderJob() cron.Job {
	handler := cr.NewAssignOrderCommandHandler()
	job, err := jobs.NewAssignOrderJob(handler)
	if err != nil {
		panic(err)
	}
	return job
}

func (cr *CompositionRoot) NewMoveCouriersJob() cron.Job {
	handler := cr.NewMoveCouriersCommandHandler()
	job, err := jobs.NewMoveCouriersJob(handler)
	if err != nil {
		panic(err)
	}
	return job
}

func (cr *CompositionRoot) NewAssignOrderCommandHandler() commands.AssignOrderCommandHandler {
	txManager := cr.newTxManager()
	orderRepository := cr.newOrderRepository(txManager)
	courierRepository := cr.newCourierRepository(txManager)

	handler, err := commands.NewAssignOrderCommandHandler(
		txManager,
		cr.newOrderDispatcher(),
		orderRepository,
		courierRepository,
	)
	if err != nil {
		panic(err)
	}
	return handler
}

func (cr *CompositionRoot) NewCreateCourierCommandHandler() commands.CreateCourierCommandHandler {
	txManager := cr.newTxManager()
	courierRepository := cr.newCourierRepository(txManager)

	handler, err := commands.NewCreateCourierCommandHandler(txManager, courierRepository)
	if err != nil {
		panic(err)
	}
	return handler
}

func (cr *CompositionRoot) NewCreateOrderCommandHandler() commands.CreateOrderCommandHandler {
	txManager := cr.newTxManager()
	orderRepository := cr.newOrderRepository(txManager)

	handler, err := commands.NewCreateOrderCommandHandler(txManager, orderRepository)
	if err != nil {
		panic(err)
	}
	return handler
}

func (cr *CompositionRoot) NewMoveCouriersCommandHandler() commands.MoveCouriersCommandHandler {
	txManager := cr.newTxManager()
	orderRepository := cr.newOrderRepository(txManager)
	courierRepository := cr.newCourierRepository(txManager)

	handler, err := commands.NewMoveCouriersCommandHandler(txManager, orderRepository, courierRepository)
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

func (cr *CompositionRoot) newTxManager() shared.TxManager {
	tx, err := shared.NewTxManager(cr.gormDb)
	if err != nil {
		panic(err)
	}
	return tx
}

func (cr *CompositionRoot) newOrderDispatcher() services.OrderDispatcher {
	return services.NewOrderDispatcher()
}

func (cr *CompositionRoot) newOrderRepository(txManager shared.TxManager) ports.OrderRepository {
	res, err := orderrepo.NewOrderRepository(txManager)
	if err != nil {
		panic(err)
	}
	return res
}

func (cr *CompositionRoot) newCourierRepository(txManager shared.TxManager) ports.CourierRepository {
	res, err := courierrepo.NewCourierRepository(txManager)
	if err != nil {
		panic(err)
	}
	return res
}
