package ports

import (
	"context"
	"delivery/internal/core/domain/model/kernel"
)

type GeoLocationGateway interface {
	DefineLocation(ctx context.Context, street string) (kernel.Location, error)
}
