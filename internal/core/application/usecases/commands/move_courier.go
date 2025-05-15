package commands

import (
	"context"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
	"errors"
)

type MoveCouriersCmd struct {
	isSet bool
}

func NewMoveCouriersCmd() (MoveCouriersCmd, error) {

	return MoveCouriersCmd{
		isSet: true,
	}, nil
}

func (c MoveCouriersCmd) IsEmpty() bool {
	return !c.isSet
}

type MoveCouriersCommandHandler interface {
	Handle(context.Context, MoveCouriersCmd) error
}

var _ MoveCouriersCommandHandler = &moveCouriersCommandHandler{}

type moveCouriersCommandHandler struct {
	unitOfWork       ports.UnitOfWork
	orderRepository  ports.OrderRepository
	courseRepository ports.CourierRepository
}

func NewMoveCouriersCommandHandler(
	unitOfWork ports.UnitOfWork,
	orderRepository ports.OrderRepository,
	courierRepository ports.CourierRepository,
) (MoveCouriersCommandHandler, error) {

	if unitOfWork == nil {
		return nil, errs.NewValueIsRequiredError("unitOfWork")
	}

	if orderRepository == nil {
		return nil, errs.NewValueIsRequiredError("orderRepository")
	}

	if courierRepository == nil {
		return nil, errs.NewValueIsRequiredError("courierRepository")
	}

	return &moveCouriersCommandHandler{
		unitOfWork:       unitOfWork,
		orderRepository:  orderRepository,
		courseRepository: courierRepository}, nil
}

func (ch *moveCouriersCommandHandler) Handle(ctx context.Context, cmd MoveCouriersCmd) error {
	if cmd.IsEmpty() {
		return errs.NewValueIsRequiredError("cmd")
	}

	assignedOrders, err := ch.orderRepository.GetAllInAssignedStatus(ctx)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return nil
		}
		return err
	}

	ch.unitOfWork.Begin(ctx)
	for _, assignedOrder := range assignedOrders {
		courier, err := ch.courseRepository.Get(ctx, *assignedOrder.CourierID())
		if err != nil {
			return err
		}

		err = courier.Move(assignedOrder.Location())
		if err != nil {
			return err
		}

		if courier.Location().Equals(assignedOrder.Location()) {
			err := assignedOrder.Complete()
			if err != nil {
				return err
			}
			err = courier.CompleteOrder(assignedOrder)
			if err != nil {
				return err
			}
		}

		err = ch.orderRepository.Update(ctx, assignedOrder)
		if err != nil {
			return err
		}
		err = ch.courseRepository.Update(ctx, courier)
		if err != nil {
			return err
		}
	}
	err = ch.unitOfWork.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
