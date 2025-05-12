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
	unitOfWork ports.UnitOfWork
}

func NewMoveCouriersCommandHandler(
	unitOfWork ports.UnitOfWork) (MoveCouriersCommandHandler, error) {
	if unitOfWork == nil {
		return nil, errs.NewValueIsRequiredError("unitOfWork")
	}

	return &moveCouriersCommandHandler{
		unitOfWork: unitOfWork}, nil
}

func (ch *moveCouriersCommandHandler) Handle(ctx context.Context, cmd MoveCouriersCmd) error {
	if cmd.IsEmpty() {
		return errs.NewValueIsRequiredError("cmd")
	}

	assignedOrders, err := ch.unitOfWork.OrderRepository().GetAllInAssignedStatus(ctx)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return nil
		}
		return err
	}

	ch.unitOfWork.Begin(ctx)
	for _, assignedOrder := range assignedOrders {
		courier, err := ch.unitOfWork.CourierRepository().Get(ctx, *assignedOrder.CourierID())
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

		err = ch.unitOfWork.OrderRepository().Update(ctx, assignedOrder)
		if err != nil {
			return err
		}
		err = ch.unitOfWork.CourierRepository().Update(ctx, courier)
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
