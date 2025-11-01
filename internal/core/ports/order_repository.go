package ports

import (
	"context"
	"delivery/internal/core/domain/model/order"
	"github.com/google/uuid"
)

type OrderRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*order.Order, error)
	GetFirstInCreatedStatus(ctx context.Context) (*order.Order, error)
	GetAllInAssignedStatus(ctx context.Context) ([]*order.Order, error)
	Save(ctx context.Context, orders ...*order.Order) error
}
