package commands

import (
	"context"
	"delivery/internal/core/domain/services"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
	"errors"
)

var (
	NotAvailableOrders   = errors.New("not available orders")
	NotAvailableCouriers = errors.New("not available couriers")
)

type AssignOrderCmd struct {
	isSet bool
}

func NewAssignOrdersCommand() AssignOrderCmd {
	return AssignOrderCmd{isSet: true}
}

func (c AssignOrderCmd) IsEmpty() bool {
	return !c.isSet
}

type AssignOrderCommandHandler interface {
	Handle(context.Context, AssignOrderCmd) error
}

var _ AssignOrderCommandHandler = &assignOrdersCommandHandler{}

type assignOrdersCommandHandler struct {
	unitOfWork      ports.UnitOfWork
	orderDispatcher services.OrderDispatcher
}

func NewAssignOrderCommandHandler(
	unitOfWork ports.UnitOfWork,
	orderDispatcher services.OrderDispatcher) (AssignOrderCommandHandler, error) {
	if unitOfWork == nil {
		return nil, errs.NewValueIsRequiredError("unitOfWork")
	}
	if orderDispatcher == nil {
		return nil, errs.NewValueIsRequiredError("orderDispatcher")
	}

	return &assignOrdersCommandHandler{
		unitOfWork:      unitOfWork,
		orderDispatcher: orderDispatcher}, nil
}

func (ch *assignOrdersCommandHandler) Handle(ctx context.Context, command AssignOrderCmd) error {
	if command.IsEmpty() {
		return errs.NewValueIsRequiredError("command")
	}

	orderAggregate, err := ch.unitOfWork.OrderRepository().GetFirstInCreatedStatus(ctx)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return NotAvailableOrders
		}
		return err
	}

	couriers, err := ch.unitOfWork.CourierRepository().GetAllFree(ctx)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return NotAvailableCouriers
		}
		return err
	}
	if couriers == nil || len(couriers) == 0 {
		return NotAvailableCouriers
	}

	courier, err := ch.orderDispatcher.Dispatch(orderAggregate, couriers)
	if err != nil {
		return err
	}

	ch.unitOfWork.Begin(ctx)

	err = ch.unitOfWork.OrderRepository().Update(ctx, orderAggregate)
	if err != nil {
		return err
	}
	err = ch.unitOfWork.CourierRepository().Update(ctx, courier)
	if err != nil {
		return err
	}

	err = ch.unitOfWork.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}
