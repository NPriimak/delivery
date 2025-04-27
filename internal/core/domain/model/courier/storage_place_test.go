package courier

import (
	"delivery/internal/pkg/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_givenValidParams_whenCreateStoragePlace_thenSuccess(t *testing.T) {
	name := "Bag"
	volume := 5

	result, err := NewStoragePlace(name, volume)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.NotEmpty(t, result.ID())
	assert.Equal(t, name, result.Name())
	assert.Equal(t, volume, result.TotalVolume())
}

func Test_givenInvalidParams_whenCreateStoragePlace_thenReturnError(t *testing.T) {
	name := "Bag"

	tests := map[string]struct {
		name     string
		volume   int
		expected error
	}{
		"name_is_blank":         {"  ", 1, errs.NewValueIsRequiredError("name")},
		"volume_is_0":           {name, 0, errs.NewValueIsRequiredError("totalVolume")},
		"volume_is_less_then_1": {name, -100, errs.NewValueIsRequiredError("totalVolume")},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			_, err := NewStoragePlace(test.name, test.volume)

			assert.Error(t, err)
			assert.Errorf(t, err, test.expected.Error())
		})
	}
}

func Test_whenAskCanStore_thenReturnCorrectAnswer(t *testing.T) {
	occupied, _ := NewStoragePlace("Bag", 5)
	empty, _ := NewStoragePlace("Bag", 5)

	_ = occupied.Store(uuid.New(), 3)

	tests := map[string]struct {
		storage  *StoragePlace
		expected bool
	}{
		"storage_is_occupied": {occupied, false},
		"storage_is_empty":    {empty, true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result, _ := test.storage.CanStore(4)

			assert.Equal(t, test.expected, result)
		})
	}
}

func Test_givenTwoStoragePlaces_whenEquals_thenReturnCorrectResult(t *testing.T) {
	place1, _ := NewStoragePlace("Place1", 5)
	place2, _ := NewStoragePlace("Place2", 5)
	place1Copy := &StoragePlace{
		id:          place1.id,
		name:        "Copy",
		totalVolume: 10,
	}

	tests := map[string]struct {
		a        *StoragePlace
		b        *StoragePlace
		expected bool
	}{
		"same_instance":          {place1, place1, true},
		"different_ids":          {place1, place2, false},
		"same_id_different_data": {place1, place1Copy, true},
		"nil_comparison":         {place1, nil, false},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := test.a.Equals(test.b)
			assert.Equal(t, test.expected, result)
		})
	}
}

func Test_givenValidOrder_whenStore_thenSuccess(t *testing.T) {
	storage, _ := NewStoragePlace("Bag", 5)
	orderID := uuid.New()
	volume := 3

	err := storage.Store(orderID, volume)

	assert.NoError(t, err)
	assert.NotNil(t, storage.OrderID())
	assert.Equal(t, orderID, *storage.OrderID())
}

func Test_givenInvalidParams_whenStore_thenReturnError(t *testing.T) {
	storage, _ := NewStoragePlace("Bag", 5)
	validOrderID := uuid.New()
	invalidOrderID := uuid.Nil
	validVolume := 3
	invalidVolume := 0
	tooBigVolume := 10

	tests := map[string]struct {
		orderID  uuid.UUID
		volume   int
		expected error
	}{
		"nil_order_id":   {invalidOrderID, validVolume, errs.NewValueIsRequiredError("orderID")},
		"invalid_volume": {validOrderID, invalidVolume, errs.NewValueIsRequiredError("volume")},
		"volume_too_big": {validOrderID, tooBigVolume, ErrCannotStoreOrderInThisStoragePlace},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := storage.Store(test.orderID, test.volume)

			assert.Error(t, err)
			assert.Errorf(t, err, test.expected.Error())
		})
	}
}

func Test_givenOccupiedStorage_whenStore_thenReturnError(t *testing.T) {
	storage, _ := NewStoragePlace("Bag", 5)
	firstOrderID := uuid.New()
	secondOrderID := uuid.New()

	err := storage.Store(firstOrderID, 3)
	assert.NoError(t, err)

	err = storage.Store(secondOrderID, 2)
	assert.Error(t, err)
	assert.Equal(t, ErrCannotStoreOrderInThisStoragePlace.Error(), err.Error())
}

func Test_givenValidOrderAndEnoughSpace_whenStore_thenSuccess(t *testing.T) {
	storage, _ := NewStoragePlace("LargeBag", 10)
	orderID := uuid.New()
	volume := 7

	err := storage.Store(orderID, volume)

	assert.NoError(t, err)
	assert.NotEmpty(t, storage.OrderID())
	assert.Equal(t, orderID, *storage.OrderID())
}

func Test_givenStoredOrder_whenClear_thenSuccess(t *testing.T) {
	storage, _ := NewStoragePlace("Bag", 5)
	orderID := uuid.New()
	_ = storage.Store(orderID, 3)

	err := storage.Clear(orderID)

	assert.NoError(t, err)
	assert.Empty(t, storage.OrderID())
}

func Test_givenInvalidParams_whenClear_thenReturnError(t *testing.T) {
	storage, _ := NewStoragePlace("Bag", 5)
	storedOrderID := uuid.New()
	wrongOrderID := uuid.New()
	_ = storage.Store(storedOrderID, 3)

	tests := map[string]struct {
		orderID  uuid.UUID
		expected error
	}{
		"nil_order_id":   {uuid.Nil, errs.NewValueIsRequiredError("orderID")},
		"wrong_order_id": {wrongOrderID, ErrOrderNotStoredInThisPlace},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := storage.Clear(test.orderID)

			assert.Error(t, err)
			assert.Errorf(t, err, test.expected.Error())
		})
	}
}
