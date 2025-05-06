package postgres

import (
	"delivery/internal/adapters/out/postgres/courierrepo"
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/order"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm/clause"
	"testing"
)

func Test_CourierRepository_Add(t *testing.T) {
	t.Run("Add new courier", func(t *testing.T) {
		ctx, db := setupTest(t)
		uow := createUOW(t, db)

		location := createTestLocation(t, 1, 1)
		expected, err := courier.NewCourier("test", 5, location)
		assert.NoError(t, err)

		err = uow.CourierRepository().Add(ctx, expected)
		assert.NoError(t, err)

		var result courierrepo.CourierDTO
		err = db.First(&result, "id = ?", expected.ID()).Error
		assert.NoError(t, err)

		assert.Equal(t, expected.ID(), result.ID)
		assert.Equal(t, expected.Name(), result.Name)
		assert.Equal(t, expected.Speed(), result.Speed)
		assert.Equal(t, expected.Location().X(), result.Location.X)
		assert.Equal(t, expected.Location().Y(), result.Location.Y)
	})

	t.Run("Add new courier with storage", func(t *testing.T) {
		ctx, db := setupTest(t)
		uow := createUOW(t, db)

		location := createTestLocation(t, 1, 1)
		expected, err := courier.NewCourier("test", 5, location)
		assert.NoError(t, err)
		err = expected.AddStoragePlace("Bag", 5)
		assert.NoError(t, err)

		err = uow.CourierRepository().Add(ctx, expected)
		assert.NoError(t, err)

		var result courierrepo.CourierDTO
		err = db.WithContext(ctx).Preload(clause.Associations).First(&result, "id = ?", expected.ID()).Error
		assert.NoError(t, err)

		assert.Equal(t, len(expected.StoragePlaces()), len(expected.StoragePlaces()))
		assert.Equal(t, len(expected.StoragePlaces()[0].ID()), len(expected.StoragePlaces()[0].ID()))
		assert.Equal(t, expected.StoragePlaces()[0].Name(), result.StoragePlaces[0].Name)
		assert.Equal(t, expected.StoragePlaces()[0].TotalVolume(), result.StoragePlaces[0].TotalVolume)
	})
}

func Test_CourierRepository_Update(t *testing.T) {
	t.Run("Update courier", func(t *testing.T) {
		ctx, db := setupTest(t)
		uow := createUOW(t, db)

		location := createTestLocation(t, 1, 1)
		old, err := courier.NewCourier("test", 5, location)
		assert.NoError(t, err)

		err = db.Create(courierrepo.DomainToDTO(old)).Error
		assert.NoError(t, err)

		err = old.AddStoragePlace("Bag", 5)
		assert.NoError(t, err)

		err = uow.CourierRepository().Update(ctx, old)
		assert.NoError(t, err)

		var result courierrepo.CourierDTO
		err = db.WithContext(ctx).Preload(clause.Associations).First(&result, "id = ?", old.ID()).Error
		assert.NoError(t, err)

		assert.Equal(t, len(old.StoragePlaces()), len(old.StoragePlaces()))
		assert.Equal(t, len(old.StoragePlaces()[0].ID()), len(old.StoragePlaces()[0].ID()))
		assert.Equal(t, old.StoragePlaces()[0].Name(), result.StoragePlaces[0].Name)
		assert.Equal(t, old.StoragePlaces()[0].TotalVolume(), result.StoragePlaces[0].TotalVolume)
	})
}

func Test_CourierRepository_GetByID(t *testing.T) {
	t.Run("Get courier by ID", func(t *testing.T) {
		ctx, db := setupTest(t)
		uow := createUOW(t, db)

		location := createTestLocation(t, 1, 1)
		expected, err := courier.NewCourier("test", 5, location)
		assert.NoError(t, err)

		err = db.Create(courierrepo.DomainToDTO(expected)).Error
		assert.NoError(t, err)

		result, err := uow.CourierRepository().Get(ctx, expected.ID())
		assert.NoError(t, err)

		assert.Equal(t, expected.ID(), result.ID())
	})
}

func Test_CourierRepository_GetAllFree(t *testing.T) {
	t.Run("Return all free couriers", func(t *testing.T) {
		ctx, db := setupTest(t)
		uow := createUOW(t, db)

		location := createTestLocation(t, 1, 1)
		first, _ := courier.NewCourier("test", 5, location)
		second, _ := courier.NewCourier("test", 6, location)

		order, _ := order.NewOrder(uuid.New(), createTestLocation(t, 10, 10), 5)
		third, _ := courier.NewCourier("test", 7, location)
		third.AddStoragePlace("Bag", 5)
		third.TakeOrder(order)

		db.Create(courierrepo.DomainToDTO(first)).
			Create(courierrepo.DomainToDTO(second)).
			Create(courierrepo.DomainToDTO(third))

		result, err := uow.CourierRepository().GetAllFree(ctx)
		assert.NoError(t, err)

		assert.Len(t, result, 2)
		notContinsBusyCourier := true
		for _, c := range result {
			if c.ID() == third.ID() {
				notContinsBusyCourier = false
			}
		}

		assert.True(t, notContinsBusyCourier)
	})
}
