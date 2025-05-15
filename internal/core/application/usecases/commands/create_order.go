package commands

import (
	"context"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
	"github.com/google/uuid"
	"strings"
)

type CreateOrderCmd struct {
	orderID uuid.UUID
	street  string
	volume  int

	isSet bool
}

func NewCreateOrderCmd(orderID uuid.UUID, street string, volume int) (CreateOrderCmd, error) {
	if orderID == uuid.Nil {
		return CreateOrderCmd{isSet: false}, errs.NewValueIsRequiredError("orderID")
	}

	if strings.TrimSpace(street) == "" {
		return CreateOrderCmd{isSet: false}, errs.NewValueIsRequiredError("street")
	}

	if volume <= 0 {
		return CreateOrderCmd{isSet: false}, errs.NewValueIsRequiredError("volume")
	}

	return CreateOrderCmd{
		orderID: orderID,
		street:  street,
		volume:  volume,
		isSet:   true,
	}, nil
}

func (cmd CreateOrderCmd) OrderID() uuid.UUID {
	return cmd.orderID
}

func (cmd CreateOrderCmd) Street() string {
	return cmd.street
}

func (cmd CreateOrderCmd) Volume() int {
	return cmd.volume
}

func (cmd CreateOrderCmd) IsEmpty() bool {
	return !cmd.isSet
}

type CreateOrderCommandHandler interface {
	Handle(ctx context.Context, cmd CreateOrderCmd) error
}

var _ CreateOrderCommandHandler = &createOrderCommandHandler{}

type createOrderCommandHandler struct {
	unitOfWork         ports.UnitOfWork
	orderRepository    ports.OrderRepository
	geoLocationGateway ports.GeoLocationGateway
}

func NewCreateOrderCommandHandler(
	uow ports.UnitOfWork,
	repo ports.OrderRepository,
	geoLocationGateway ports.GeoLocationGateway,
) (CreateOrderCommandHandler, error) {
	if uow == nil {
		return nil, errs.NewValueIsRequiredError("uow")
	}

	if repo == nil {
		return nil, errs.NewValueIsRequiredError("repo")
	}

	if geoLocationGateway == nil {
		return nil, errs.NewValueIsRequiredError("geoLocationGateway")
	}
	return &createOrderCommandHandler{
		unitOfWork:         uow,
		orderRepository:    repo,
		geoLocationGateway: geoLocationGateway,
	}, nil
}

func (ch *createOrderCommandHandler) Handle(ctx context.Context, cmd CreateOrderCmd) error {
	if cmd.IsEmpty() {
		return errs.NewValueIsRequiredError("cmd")
	}

	existingOrder, err := ch.orderRepository.Get(ctx, cmd.OrderID())
	if err != nil {
		return err
	}
	if existingOrder != nil {
		return nil
	}

	location, err := ch.geoLocationGateway.DefineLocation(ctx, cmd.Street())
	if err != nil {
		return err
	}

	existingOrder, err = order.NewOrder(
		cmd.OrderID(),
		location,
		cmd.Volume(),
	)
	if err != nil {
		return err
	}

	err = ch.orderRepository.Add(ctx, existingOrder)
	if err != nil {
		return err
	}

	return nil
}
