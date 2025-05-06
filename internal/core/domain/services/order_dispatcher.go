package services

import (
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"
	"errors"
	"math"
)

var (
	SuitableCourierNotFound = errors.New("suitable courier not found")
)

type OrderDispatcher interface {
	Dispatch(order *order.Order, couriers []*courier.Courier) (*courier.Courier, error)
}

var _ OrderDispatcher = &orderDispatcher{}

type orderDispatcher struct {
}

func NewOrderDispatcher() OrderDispatcher {
	return &orderDispatcher{}
}

func (p *orderDispatcher) Dispatch(currentOrder *order.Order, couriers []*courier.Courier) (*courier.Courier, error) {
	if currentOrder == nil {
		return nil, errs.NewValueIsRequiredError("currentOrder")
	}
	if couriers == nil || len(couriers) == 0 {
		return nil, errs.NewValueIsRequiredError("couriers")
	}

	bestCourier, err := p.findBestCourier(currentOrder, couriers)
	if err != nil {
		return nil, err
	}

	if currentOrder.Status() != order.StatusCreated {
		return nil, order.ErrOrderHasAlreadyBeenAssigned
	}
	if err := currentOrder.Assign(bestCourier.ID()); err != nil {
		return nil, err
	}
	if err := bestCourier.TakeOrder(currentOrder); err != nil {
		return nil, err
	}

	return bestCourier, nil
}

func (p *orderDispatcher) findBestCourier(order *order.Order, couriers []*courier.Courier) (*courier.Courier, error) {
	var bestCourier *courier.Courier
	minTime := math.MaxFloat64

	for _, c := range couriers {
		canTake, err := c.CanTakeOrder(order)
		if err != nil {
			return nil, err
		}
		if !canTake {
			continue
		}

		time, err := c.CalculateTimeToLocation(order.Location())
		if err != nil {
			return nil, err
		}

		if time < minTime {
			minTime = time
			bestCourier = c
		}
	}

	if bestCourier == nil {
		return nil, SuitableCourierNotFound
	}

	return bestCourier, nil
}
