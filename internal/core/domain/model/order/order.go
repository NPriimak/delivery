package order

import (
	"delivery/internal/core/domain/kernel"
	"delivery/internal/pkg/errs"
	"errors"
	"github.com/google/uuid"
)

var (
	ErrOrderHasAlreadyBeenAssigned  = errors.New("order has already been assigned")
	ErrOrderHasNotBeenAssigned      = errors.New("order has not been assigned")
	ErrOrderHasAlreadyBeenCompleted = errors.New("Order has already been completed")
)

type Order struct {
	id        uuid.UUID
	courierID *uuid.UUID
	location  kernel.Location
	volume    int
	status    Status
}

func NewOrder(orderID uuid.UUID, location kernel.Location, volume int) (*Order, error) {
	if orderID == uuid.Nil {
		return nil, errs.NewValueIsRequiredError("orderID")
	}
	if location.IsEmpty() {
		return nil, errs.NewValueIsRequiredError("location")
	}
	if volume <= 0 {
		return nil, errs.NewValueIsRequiredError("volume")
	}

	return &Order{
		id:       orderID,
		location: location,
		volume:   volume,
		status:   StatusCreated,
	}, nil
}

func (o *Order) Assign(courierID uuid.UUID) error {
	if courierID == uuid.Nil {
		return errs.NewValueIsRequiredError("courierID")
	}

	if o.status != StatusCreated {
		return ErrOrderHasAlreadyBeenAssigned
	}

	o.courierID = &courierID
	o.status = StatusAssigned
	return nil
}

func (o *Order) Complete() error {
	if !o.isAssigned() {
		return ErrOrderHasNotBeenAssigned
	}

	if o.status == StatusCompleted {
		return ErrOrderHasAlreadyBeenCompleted
	}

	o.status = StatusCompleted
	return nil
}

func (o *Order) isAssigned() bool {
	return o.courierID != nil && o.status == StatusAssigned
}

func (o *Order) ID() uuid.UUID {
	return o.id
}

func (o *Order) CourierID() *uuid.UUID {
	return o.courierID
}

func (o *Order) Location() kernel.Location {
	return o.location
}

func (o *Order) Volume() int {
	return o.volume
}

func (o *Order) Status() Status {
	return o.status
}

func (o *Order) Equals(other *Order) bool {
	if other == nil {
		return false
	}

	return o.id == other.id
}
