package ports

import (
	"context"
)

type UnitOfWork interface {
	Begin(ctx context.Context)
	Commit(ctx context.Context) error

	CourierRepository() CourierRepository
	OrderRepository() OrderRepository
}
