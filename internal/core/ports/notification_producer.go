package ports

import (
	"context"
	"delivery/internal/pkg/ddd"
)

type NotificationProducer interface {
	Publish(ctx context.Context, domainEvent ddd.DomainEvent) error
	Close() error
}
