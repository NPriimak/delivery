package postgres

import (
	"delivery/internal/adapters/out/postgres/orderrepo"
	"delivery/internal/adapters/out/postgres/shared"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/core/ports"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_OrderRepository_Add(t *testing.T) {
	t.Run("Must add new order", func(t *testing.T) {
		ctx, db := setupTest(t)
		tx := createTxManager(t, db)
		repo := createOrderRepository(t, tx)

		newOrder, err := order.NewOrder(uuid.New(), createTestLocation(t, 1, 1), 5)
		assert.NoError(t, err)

		err = repo.Add(ctx, newOrder)
		assert.NoError(t, err)

		var orderFromDb orderrepo.OrderDTO
		err = db.First(&orderFromDb, "id = ?", newOrder.ID()).Error
		assert.NoError(t, err)

		assert.Equal(t, newOrder.ID(), orderFromDb.ID)
		assert.Equal(t, newOrder.Volume(), orderFromDb.Volume)
		assert.Equal(t, newOrder.Status(), orderFromDb.Status)
		assert.Equal(t, newOrder.Location().X(), orderFromDb.Location.X)
		assert.Equal(t, newOrder.Location().Y(), orderFromDb.Location.Y)
	})
}

func Test_OrderRepository_Update(t *testing.T) {
	t.Run("Must update order", func(t *testing.T) {
		ctx, db := setupTest(t)
		tx := createTxManager(t, db)
		repo := createOrderRepository(t, tx)

		location := createTestLocation(t, 1, 1)
		oldOrder, err := order.NewOrder(uuid.New(), location, 5)
		assert.NoError(t, err)
		dto := orderrepo.DomainToDTO(oldOrder)

		db.Create(&dto)

		courierID := uuid.New()
		err = oldOrder.Assign(courierID)
		assert.NoError(t, err)

		err = repo.Update(ctx, oldOrder)
		assert.NoError(t, err)

		var orderFromDb orderrepo.OrderDTO
		err = db.First(&orderFromDb, "id = ?", oldOrder.ID()).Error
		assert.NoError(t, err)

		assert.Equal(t, oldOrder.CourierID(), orderFromDb.CourierID)
		assert.Equal(t, order.StatusAssigned, orderFromDb.Status)
	})
}

func Test_OrderRepository_GetFirstInCreatedStatus(t *testing.T) {
	t.Run("Return first in created status", func(t *testing.T) {
		ctx, db := setupTest(t)
		tx := createTxManager(t, db)
		repo := createOrderRepository(t, tx)

		location := createTestLocation(t, 1, 1)
		inCreated, err := order.NewOrder(uuid.New(), location, 5)
		assert.NoError(t, err)

		assigned, err := order.NewOrder(uuid.New(), location, 5)
		assert.NoError(t, err)
		err = assigned.Assign(uuid.New())
		assert.NoError(t, err)

		db.Create(orderrepo.DomainToDTO(inCreated))
		db.Create(orderrepo.DomainToDTO(assigned))

		result, err := repo.GetFirstInCreatedStatus(ctx)
		assert.NoError(t, err)

		assert.Equal(t, result.ID(), inCreated.ID())
	})
}

func Test_OrderRepository_GetAllAssignedOrders(t *testing.T) {
	t.Run("Return all assigned orders", func(t *testing.T) {
		ctx, db := setupTest(t)
		tx := createTxManager(t, db)
		repo := createOrderRepository(t, tx)

		location := createTestLocation(t, 1, 1)
		first, err := order.NewOrder(uuid.New(), location, 5)
		assert.NoError(t, err)
		second, err := order.NewOrder(uuid.New(), location, 5)
		assert.NoError(t, err)
		err = first.Assign(uuid.New())
		assert.NoError(t, err)
		err = second.Assign(uuid.New())
		assert.NoError(t, err)

		db.Create(orderrepo.DomainToDTO(first))
		db.Create(orderrepo.DomainToDTO(second))

		result, err := repo.GetAllInAssignedStatus(ctx)
		assert.NoError(t, err)

		assert.Len(t, result, 2)
	})
}

func Test_OrderRepository_GetByID(t *testing.T) {
	t.Run("Return order by id", func(t *testing.T) {
		ctx, db := setupTest(t)
		tx := createTxManager(t, db)
		repo := createOrderRepository(t, tx)

		location := createTestLocation(t, 1, 1)
		expected, err := order.NewOrder(uuid.New(), location, 5)
		assert.NoError(t, err)
		db.Create(orderrepo.DomainToDTO(expected))

		fromDb, err := repo.Get(ctx, expected.ID())
		assert.NoError(t, err)

		assert.Equal(t, expected.ID(), fromDb.ID())
	})
}

func createOrderRepository(t *testing.T, tx shared.TxManager) ports.OrderRepository {
	res, err := orderrepo.NewOrderRepository(tx)
	assert.NoError(t, err)
	return res
}
