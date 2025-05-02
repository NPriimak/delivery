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

type IOrderDispatcher interface {
	Dispatch(order *order.Order, couriers []*courier.Courier) (*courier.Courier, error)
}

var _ IOrderDispatcher = &OrderDispatcher{}

type OrderDispatcher struct {
}

func NewOrderDispatcher() IOrderDispatcher {
	return &OrderDispatcher{}
}

func (p *OrderDispatcher) Dispatch(order *order.Order, couriers []*courier.Courier) (*courier.Courier, error) {
	if order == nil {
		return nil, errs.NewValueIsRequiredError("order")
	}
	if couriers == nil || len(couriers) == 0 {
		return nil, errs.NewValueIsRequiredError("couriers")
	}

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

	if err := order.Assign(bestCourier.ID()); err != nil {
		return nil, err
	}
	if err := bestCourier.TakeOrder(order); err != nil {
		return nil, err
	}

	return bestCourier, nil
}
