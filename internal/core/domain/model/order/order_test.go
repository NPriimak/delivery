package order

import (
	"delivery/internal/core/domain/kernel"
	"delivery/internal/pkg/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_createOrder(t *testing.T) {
	t.Run("given valid parameters when NewOrder then success", func(t *testing.T) {
		orderID := uuid.New()
		location := createTestLocation(t)
		volume := 5

		order, err := NewOrder(orderID, location, volume)

		assert.NoError(t, err)
		assert.Equal(t, orderID, order.ID())
		assert.Equal(t, location, order.Location())
		assert.Equal(t, volume, order.Volume())
		assert.Equal(t, StatusCreated, order.Status())
		assert.Empty(t, order.CourierID())
	})
}

func Test_givenInvalidParams_whenCreateNewOrder_thenFail(t *testing.T) {
	t.Run("given invalid parameters when NewOrder then return error", func(t *testing.T) {
		validID := uuid.New()
		validLocation := createTestLocation(t)
		validVolume := 5

		tests := map[string]struct {
			id       uuid.UUID
			location kernel.Location
			volume   int
			expected error
		}{
			"nil_order_id":    {uuid.Nil, validLocation, validVolume, errs.NewValueIsRequiredError("orderID")},
			"empty_location":  {validID, kernel.Location{}, validVolume, errs.NewValueIsRequiredError("location")},
			"zero_volume":     {validID, validLocation, 0, errs.NewValueIsRequiredError("volume")},
			"negative_volume": {validID, validLocation, -1, errs.NewValueIsRequiredError("volume")},
		}

		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				_, err := NewOrder(test.id, test.location, test.volume)
				assert.Error(t, err)
				assert.Errorf(t, err, test.expected.Error())
			})
		}
	})
}

func Test_assignOrder(t *testing.T) {
	t.Run("given unassigned order when Assign then success", func(t *testing.T) {
		order := createTestOrder(t)
		courierID := uuid.New()

		err := order.Assign(courierID)

		assert.NoError(t, err)
		assert.Equal(t, courierID, *order.CourierID())
		assert.Equal(t, StatusAssigned, order.Status())
	})

	t.Run("given invalid courierID when Assign then return error", func(t *testing.T) {
		order := createTestOrder(t)

		err := order.Assign(uuid.Nil)
		assert.Error(t, err)
		assert.Errorf(t, err, errs.NewValueIsRequiredError("courierID").Error())
	})

	t.Run("given already assigned order when Assign then return error", func(t *testing.T) {
		order := createTestOrder(t)
		courierID := uuid.New()
		_ = order.Assign(courierID)

		err := order.Assign(uuid.New())
		assert.Error(t, err)
		assert.Errorf(t, err, ErrOrderHasAlreadyBeenAssigned.Error())
	})
}

func Test_completeOrder(t *testing.T) {
	t.Run("given assigned order when Complete then success", func(t *testing.T) {
		order := createTestOrder(t)
		_ = order.Assign(uuid.New())

		err := order.Complete()
		assert.NoError(t, err)
		assert.Equal(t, StatusCompleted, order.Status())
	})

	t.Run("given unassigned order when Complete then return error", func(t *testing.T) {
		order := createTestOrder(t)

		err := order.Complete()
		assert.Error(t, err)
		assert.Errorf(t, err, ErrOrderHasNotBeenAssigned.Error())
		assert.Equal(t, StatusCreated, order.Status())
	})

	t.Run("given order has already been completed", func(t *testing.T) {
		order := createTestOrder(t)
		order.status = StatusCompleted

		err := order.Complete()
		assert.Error(t, err)
		assert.Errorf(t, err, ErrOrderHasAlreadyBeenCompleted.Error())
	})
}

func Test_equals(t *testing.T) {
	t.Run("given two orders when Equals then return correct result", func(t *testing.T) {
		order1 := createTestOrder(t)
		order2 := createTestOrder(t)
		order1Copy := &Order{
			id:        order1.id,
			courierID: nil,
			location:  createTestLocation(t),
			volume:    100,
			status:    StatusCompleted,
		}

		tests := map[string]struct {
			a        *Order
			b        *Order
			expected bool
		}{
			"same_instance":          {order1, order1, true},
			"different_ids":          {order1, order2, false},
			"same_id_different_data": {order1, order1Copy, true},
			"nil_comparison":         {order1, nil, false},
		}

		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				result := test.a.Equals(test.b)
				assert.Equal(t, test.expected, result)
			})
		}
	})
}

func createTestOrder(t *testing.T) *Order {
	location := createTestLocation(t)
	order, err := NewOrder(uuid.New(), location, 5)
	assert.NoError(t, err)
	return order
}

func createTestLocation(t *testing.T) kernel.Location {
	location, err := kernel.NewLocation(7, 10)
	assert.NoError(t, err)
	return location
}
