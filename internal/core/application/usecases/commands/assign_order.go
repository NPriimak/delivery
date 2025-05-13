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
	unitOfWork        ports.UnitOfWork
	orderRepository   ports.OrderRepository
	courierRepository ports.CourierRepository
	orderDispatcher   services.OrderDispatcher
}

func NewAssignOrderCommandHandler(
	unitOfWork ports.UnitOfWork,
	orderDispatcher services.OrderDispatcher,
	orderRepository ports.OrderRepository,
	courierRepository ports.CourierRepository,
) (AssignOrderCommandHandler, error) {
	if unitOfWork == nil {
		return nil, errs.NewValueIsRequiredError("unitOfWork")
	}
	if orderRepository == nil {
		return nil, errs.NewValueIsRequiredError("orderRepository")
	}
	if courierRepository == nil {
		return nil, errs.NewValueIsRequiredError("courierRepository")
	}
	if orderDispatcher == nil {
		return nil, errs.NewValueIsRequiredError("orderDispatcher")
	}

	return &assignOrdersCommandHandler{
		unitOfWork:        unitOfWork,
		orderRepository:   orderRepository,
		courierRepository: courierRepository,
		orderDispatcher:   orderDispatcher}, nil
}

func (ch *assignOrdersCommandHandler) Handle(ctx context.Context, command AssignOrderCmd) error {
	if command.IsEmpty() {
		return errs.NewValueIsRequiredError("command")
	}

	orderAggregate, err := ch.orderRepository.GetFirstInCreatedStatus(ctx)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return NotAvailableOrders
		}
		return err
	}

	couriers, err := ch.courierRepository.GetAllFree(ctx)
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

	err = ch.orderRepository.Update(ctx, orderAggregate)
	if err != nil {
		return err
	}
	err = ch.courierRepository.Update(ctx, courier)
	if err != nil {
		return err
	}

	err = ch.unitOfWork.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}
