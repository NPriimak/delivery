package commands

import (
	"context"
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"
	"delivery/mocks/core/ports"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_CreateNewHandler(t *testing.T) {
	t.Run("Return err with nil UOW", func(t *testing.T) {
		_, err := NewMoveCouriersCommandHandler(nil)
		assert.Errorf(t, err, errs.NewValueIsRequiredError("unitOfWork").Error())
	})

}

func Test_Handle_WithInvalidArgs(t *testing.T) {
	t.Run("return error if given empty cmd", func(t *testing.T) {
		uow := ports.NewMockUnitOfWork(t)
		handler, _ := NewMoveCouriersCommandHandler(uow)

		err := handler.Handle(context.Background(), MoveCouriersCmd{})
		assert.Errorf(t, err, errs.NewValueIsRequiredError("cmd").Error())
	})
}

func Test_Handle_NegativeScenarios(t *testing.T) {
	t.Run("If order not found - return nil", func(t *testing.T) {
		uow := ports.NewMockUnitOfWork(t)
		orderRepo := ports.NewMockOrderRepository(t)

		uow.On("OrderRepository").Return(orderRepo)

		orderRepo.On("GetAllInAssignedStatus", mock.Anything).Return([]*order.Order{}, errs.ErrObjectNotFound)

		handler, _ := NewMoveCouriersCommandHandler(uow)
		cmd, _ := NewMoveCouriersCmd()

		err := handler.Handle(context.Background(), cmd)
		assert.NoError(t, err)
	})

	t.Run("If courier not found - return err", func(t *testing.T) {
		orderId := uuid.New()
		courierId := uuid.New()

		testOrder := createAssignedTestOrder(t, orderId, courierId, 1)

		uow := ports.NewMockUnitOfWork(t)
		orderRepo := ports.NewMockOrderRepository(t)
		courierRepo := ports.NewMockCourierRepository(t)

		uow.On("OrderRepository").Return(orderRepo)
		uow.On("CourierRepository").Return(courierRepo)
		uow.On("Begin", mock.Anything).Return()

		orderRepo.On("GetAllInAssignedStatus", mock.Anything).Return([]*order.Order{testOrder}, nil)
		courierRepo.On("Get", mock.Anything, courierId).Return(nil, errs.ErrObjectNotFound)

		handler, _ := NewMoveCouriersCommandHandler(uow)
		cmd, _ := NewMoveCouriersCmd()

		err := handler.Handle(context.Background(), cmd)
		assert.ErrorIs(t, err, errs.ErrObjectNotFound)
	})
}

func Test_Handle_PositiveScenarios(t *testing.T) {
	t.Run("Move courier if ok", func(t *testing.T) {
		orderId := uuid.New()

		testOrder := createTestOrder(t, orderId, 5, createTestLocation(t, 10, 10))
		testCourier := createTestCourier(t, createTestLocation(t, 1, 1), 5)
		_ = testOrder.Assign(testCourier.ID())

		uow := ports.NewMockUnitOfWork(t)
		orderRepo := ports.NewMockOrderRepository(t)
		courierRepo := ports.NewMockCourierRepository(t)

		uow.On("OrderRepository").Return(orderRepo)
		uow.On("CourierRepository").Return(courierRepo)
		uow.On("Begin", mock.Anything).Return()
		uow.On("Commit", mock.Anything).Return(nil)

		orderRepo.On("GetAllInAssignedStatus", mock.Anything).Return([]*order.Order{testOrder}, nil)
		courierRepo.On("Get", mock.Anything, testCourier.ID()).Return(testCourier, nil)

		orderRepo.On("Update", mock.Anything, testOrder).Return(nil)
		courierRepo.On("Update", mock.Anything, testCourier).Return(nil)

		handler, _ := NewMoveCouriersCommandHandler(uow)
		cmd, _ := NewMoveCouriersCmd()

		err := handler.Handle(context.Background(), cmd)
		assert.NoError(t, err)

		assert.Equal(t, uint8(6), testCourier.Location().X())
		assert.Equal(t, uint8(1), testCourier.Location().Y())
	})
}

func createAssignedTestOrder(
	t *testing.T,
	id uuid.UUID,
	courierID uuid.UUID,
	volume int,
) *order.Order {
	location := createTestLocation(t, 5, 5)
	order, err := order.NewOrder(id, location, volume)
	assert.NoError(t, err)

	err = order.Assign(courierID)
	assert.NoError(t, err)

	return order
}

func createTestCourier(t *testing.T, location kernel.Location, speed int) *courier.Courier {
	res, err := courier.NewCourier("test", speed, location)
	assert.NoError(t, err)
	return res
}

func createTestOrder(t *testing.T, id uuid.UUID, volume int, location kernel.Location) *order.Order {
	res, err := order.NewOrder(id, location, volume)
	assert.NoError(t, err)
	return res
}

func createTestLocation(t *testing.T, x int, y int) kernel.Location {
	location, err := kernel.NewLocation(uint8(x), uint8(y))
	assert.NoError(t, err)
	return location
}
