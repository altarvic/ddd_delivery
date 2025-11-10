package ports

import (
	"context"
	"delivery/internal/core/domain/kernel"
)

type GeoClient interface {
	GetGeolocation(ctx context.Context, address string) (kernel.Location, error)
	Close() error
}
