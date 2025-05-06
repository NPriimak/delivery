package services_test

import (
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/core/domain/services"
	"delivery/internal/pkg/errs"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_DispatchErrors(t *testing.T) {
	svc := services.NewOrderDispatcher()

	t.Run("nil order", func(t *testing.T) {
		_, err := svc.Dispatch(nil, []*courier.Courier{})
		assert.EqualError(t, err, errs.NewValueIsRequiredError("currentOrder").Error())
	})

	t.Run("nil couriers slice", func(t *testing.T) {
		o := createOrder(t, 1, createLoc(t, 1, 1))
		_, err := svc.Dispatch(o, nil)
		assert.EqualError(t, err, errs.NewValueIsRequiredError("couriers").Error())
	})

	t.Run("empty couriers slice", func(t *testing.T) {
		o := createOrder(t, 1, createLoc(t, 1, 1))
		_, err := svc.Dispatch(o, []*courier.Courier{})
		assert.EqualError(t, err, errs.NewValueIsRequiredError("couriers").Error())
	})

	t.Run("No suitable courier", func(t *testing.T) {
		c := createCourier(t, 5, createLoc(t, 5, 5), 2)
		o := createOrder(t, 3, createLoc(t, 6, 6))

		_, err := svc.Dispatch(o, []*courier.Courier{c})
		assert.EqualError(t, err, services.SuitableCourierNotFound.Error())
	})
}

func TestDispatch_Success(t *testing.T) {
	svc := services.NewOrderDispatcher()

	t.Run("Choose right courier", func(t *testing.T) {
		tests := map[string]struct {
			order    *order.Order
			couriers []*courier.Courier
			expected int
		}{
			"Choose first courier": {
				createOrder(t, 5, createLoc(t, 10, 10)),
				[]*courier.Courier{
					createCourier(t, 5, createLoc(t, 9, 8), 5),
					createCourier(t, 5, createLoc(t, 5, 5), 5),
					createCourier(t, 5, createLoc(t, 4, 1), 5),
				},
				0,
			},
			"Choose second courier": {
				createOrder(t, 5, createLoc(t, 10, 10)),
				[]*courier.Courier{
					createCourier(t, 5, createLoc(t, 4, 6), 5),
					createCourier(t, 10, createLoc(t, 5, 5), 5),
					createCourier(t, 5, createLoc(t, 4, 1), 5),
				},
				1,
			},
			"Choose first suitable if we have two same couriers in the same location": {
				createOrder(t, 5, createLoc(t, 10, 10)),
				[]*courier.Courier{
					createCourier(t, 5, createLoc(t, 4, 6), 5),
					createCourier(t, 10, createLoc(t, 5, 5), 5),
					createCourier(t, 10, createLoc(t, 5, 5), 5),
				},
				1,
			},
		}

		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				anotherOrder := createOrder(t, 5, createLoc(t, 10, 10))
				expectedCourier := test.couriers[test.expected]
				c, _ := svc.Dispatch(test.order, test.couriers)

				assert.Equal(t, c, expectedCourier)

				canTake, _ := expectedCourier.CanTakeOrder(anotherOrder)
				assert.False(t, canTake)

				assert.Equal(t, *test.order.CourierID(), expectedCourier.ID())
			})
		}
	})
}

// HELPERS
func createLoc(t *testing.T, x, y uint8) kernel.Location {
	loc, err := kernel.NewLocation(x, y)
	if err != nil {
		t.Fatalf("failed to create location: %v", err)
	}
	return loc
}

func createOrder(t *testing.T, volume int, loc kernel.Location) *order.Order {
	o, err := order.NewOrder(uuid.New(), loc, volume)
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}
	return o
}

func createCourier(t *testing.T, speed int, loc kernel.Location, storageVolume int) *courier.Courier {
	c, err := courier.NewCourier("Test", speed, loc)
	if err != nil {
		t.Fatalf("failed to create courier: %v", err)
	}

	if err := c.AddStoragePlace("TestBag", storageVolume); err != nil {
		t.Fatalf("failed to add storage: %v", err)
	}
	return c
}
