package ports

import (
	"context"
	"delivery/internal/core/domain/model/courier"
	"github.com/google/uuid"
)

type CourierRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*courier.Courier, error)
	GetAllFree(ctx context.Context) ([]*courier.Courier, error)
	Save(ctx context.Context, couriers ...*courier.Courier) error
}
