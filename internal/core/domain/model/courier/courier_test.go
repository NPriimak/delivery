package courier

import (
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_createNewCourier(t *testing.T) {
	t.Run("given valid parameters when create courier then success", func(t *testing.T) {
		name := "Test"
		speed := 5
		loc := createTestLocation(t)

		c, err := NewCourier(name, speed, loc)

		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, c.ID())
		assert.Equal(t, name, c.Name())
		assert.Equal(t, speed, c.Speed())
		assert.Equal(t, loc, c.Location())
		assert.NotNil(t, c.StoragePlaces())
		assert.Empty(t, c.StoragePlaces())
		assert.NotNil(t, c.BaseAggregate)
	})

	t.Run("given invalid parameters when create new courier then return error", func(t *testing.T) {
		validName := "Test Courier"
		validSpeed := 5
		validLoc := createTestLocation(t)

		tests := map[string]struct {
			name     string
			speed    int
			location kernel.Location
			expected error
		}{
			"empty_name":     {"", validSpeed, validLoc, errs.NewValueIsRequiredError("name")},
			"blank_name":     {"   ", validSpeed, validLoc, errs.NewValueIsRequiredError("name")},
			"zero_speed":     {validName, 0, validLoc, errs.NewValueIsRequiredError("speed")},
			"negative_speed": {validName, -1, validLoc, errs.NewValueIsRequiredError("speed")},
			"empty_location": {validName, validSpeed, kernel.Location{}, errs.NewValueIsRequiredError("location")},
		}

		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				_, err := NewCourier(test.name, test.speed, test.location)
				assert.Errorf(t, err, test.expected.Error())
			})
		}
	})
}

func Test_addStoragePlace(t *testing.T) {
	t.Run("given valid storage place when add to courier then success", func(t *testing.T) {
		c := createTestCourier(t)
		storageName := "Bag"
		volume := 5

		err := c.AddStoragePlace(storageName, volume)

		assert.NoError(t, err)
		assert.Len(t, c.StoragePlaces(), 1)
		assert.Equal(t, storageName, c.StoragePlaces()[0].Name())
		assert.Equal(t, volume, c.StoragePlaces()[0].TotalVolume())
	})

	t.Run("given invalid storage place when add it to courier then return error", func(t *testing.T) {
		c := createTestCourier(t)
		validName := "Bag"
		validVolume := 5

		tests := map[string]struct {
			name     string
			volume   int
			expected error
		}{
			"empty_name":      {"", validVolume, errs.NewValueIsRequiredError("name")},
			"blank_name":      {"   ", validVolume, errs.NewValueIsRequiredError("name")},
			"zero_volume":     {validName, 0, errs.NewValueIsRequiredError("volume")},
			"negative_volume": {validName, -1, errs.NewValueIsRequiredError("volume")},
		}

		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				err := c.AddStoragePlace(test.name, test.volume)
				assert.Errorf(t, err, test.expected.Error())
				assert.Empty(t, c.StoragePlaces())
			})
		}
	})
}

func Test_canTakeOrder(t *testing.T) {
	t.Run("given no storage places when check can take order then return false", func(t *testing.T) {
		c := createTestCourier(t)
		o := createTestOrder(t)

		canTake, err := c.CanTakeOrder(o)

		assert.NoError(t, err)
		assert.False(t, canTake)
	})

	t.Run("given order with no space when check can take order then return false", func(t *testing.T) {
		c := createTestCourier(t)
		_ = c.AddStoragePlace("Small Bag", 3)
		o := createTestOrderWithVolume(t, 5)

		canTake, err := c.CanTakeOrder(o)

		assert.NoError(t, err)
		assert.False(t, canTake)
	})

	t.Run("given order with available space when check can take order then return true", func(t *testing.T) {
		c := createTestCourier(t)
		_ = c.AddStoragePlace("Large Bag", 10)
		o := createTestOrderWithVolume(t, 5)

		canTake, err := c.CanTakeOrder(o)

		assert.NoError(t, err)
		assert.True(t, canTake)
	})
}

func TestCourier_TakeOrder(t *testing.T) {
	t.Run("given valid order when take order then success", func(t *testing.T) {
		c := createTestCourier(t)
		_ = c.AddStoragePlace("Bag", 10)
		o := createTestOrderWithVolume(t, 5)

		err := c.TakeOrder(o)

		assert.NoError(t, err)
		assert.True(t, c.StoragePlaces()[0].isOccupied())
		assert.Equal(t, o.ID(), *c.StoragePlaces()[0].OrderID())
	})

	t.Run("given no storage places when TakeOrder then return error", func(t *testing.T) {
		c := createTestCourier(t)
		o := createTestOrder(t)

		err := c.TakeOrder(o)
		assert.Errorf(t, err, ErrNoStoragePlace.Error())
	})
}

func Test_calculateTimeToLocation(t *testing.T) {
	t.Run("given valid target when calculate time to location then return correct value", func(t *testing.T) {
		startLoc := createLocation(t, 1, 1)
		c, _ := NewCourier("test", 2, startLoc)
		c.location = startLoc
		targetLoc := createLocation(t, 4, 5)

		time, err := c.CalculateTimeToLocation(targetLoc)

		assert.NoError(t, err)
		assert.Equal(t, 3.5, time)
	})
}

func Test_moveToTargetLocation(t *testing.T) {
	tests := map[string]struct {
		startLocation    kernel.Location
		targetLocation   kernel.Location
		expectedLocation kernel.Location
		speed            int
	}{
		"move_right":      {createLocation(t, 1, 1), createLocation(t, 10, 1), createLocation(t, 6, 1), 5},
		"move_left":       {createLocation(t, 10, 1), createLocation(t, 1, 1), createLocation(t, 5, 1), 5},
		"move_right_down": {createLocation(t, 1, 1), createLocation(t, 5, 6), createLocation(t, 5, 2), 5},
		"move_left_up":    {createLocation(t, 10, 10), createLocation(t, 3, 6), createLocation(t, 5, 10), 5},
	}

	for name, test := range tests {
		t.Run("when "+name+" then reach expected location", func(t *testing.T) {
			c, _ := NewCourier("test", test.speed, test.startLocation)
			err := c.Move(test.targetLocation)

			assert.NoError(t, err)
			assert.Equal(t, test.expectedLocation, c.Location())
		})
	}
}

func Test_equals(t *testing.T) {
	t.Run("given two couriers when check equality then return correct result", func(t *testing.T) {
		c1 := createTestCourier(t)
		c2 := createTestCourier(t)
		c1Copy := &Courier{
			id:            c1.id,
			name:          "Copy",
			speed:         100,
			location:      createTestLocation(t),
			storagePlaces: []*StoragePlace{},
		}

		tests := map[string]struct {
			a        *Courier
			b        *Courier
			expected bool
		}{
			"same_instance":          {c1, c1, true},
			"different_ids":          {c1, c2, false},
			"same_id_different_data": {c1, c1Copy, true},
			"nil_comparison":         {c1, nil, false},
		}

		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				result := test.a.Equals(test.b)
				assert.Equal(t, test.expected, result)
			})
		}
	})
}

func Test_RestoreCourier(t *testing.T) {
	t.Run("Must correctly restore aggregate", func(t *testing.T) {
		expectedID := uuid.New()
		expectedName := "Name"
		expectedSpeed := 5
		expectedLocation := createLocation(t, 1, 1)
		expectedSP := make([]*StoragePlace, 0)
		for i := 1; i <= 5; i++ {
			sp, _ := NewStoragePlace(string(rune(i)), i)
			expectedSP = append(expectedSP, sp)
		}

		result := RestoreCourier(
			expectedID,
			expectedName,
			expectedSpeed,
			expectedLocation,
			expectedSP,
		)

		assert.Equal(t, result.ID(), expectedID)
		assert.Equal(t, result.Name(), expectedName)
		assert.Equal(t, result.Speed(), expectedSpeed)
		assert.Equal(t, result.Location(), expectedLocation)
		assert.Equal(t, len(result.StoragePlaces()), len(expectedSP))
	})
}

func createTestOrder(t *testing.T) *order.Order {
	o, err := order.NewOrder(uuid.New(), createTestLocation(t), 5)
	assert.NoError(t, err)
	return o
}

func createTestOrderWithVolume(t *testing.T, volume int) *order.Order {
	o, err := order.NewOrder(uuid.New(), createTestLocation(t), volume)
	assert.NoError(t, err)
	return o
}

func createTestCourier(t *testing.T) *Courier {
	c, err := NewCourier("Test Courier", 5, createLocation(t, 5, 5))
	assert.NoError(t, err)
	return c
}

func createTestLocation(t *testing.T) kernel.Location {
	loc, err := kernel.NewLocation(5, 5) // центральное положение
	assert.NoError(t, err)
	return loc
}

func createLocation(t *testing.T, x, y uint8) kernel.Location {
	loc, err := kernel.NewLocation(x, y)
	assert.NoError(t, err)
	return loc
}
