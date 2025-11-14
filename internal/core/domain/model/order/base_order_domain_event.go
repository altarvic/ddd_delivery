package order

import (
	"delivery/internal/pkg/ddd"
	"github.com/google/uuid"
	"time"
)

var _ ddd.DomainEvent = &baseOrderDomainEvent{}

type baseOrderDomainEvent struct {
	id         uuid.UUID
	name       string
	occurredAt time.Time

	orderId uuid.UUID
}

func newBaseOrderDomainEvent(name string, orderId uuid.UUID) baseOrderDomainEvent {
	return baseOrderDomainEvent{
		id:         uuid.New(),
		occurredAt: time.Now(),
		name:       name,
		orderId:    orderId,
	}
}

func (b *baseOrderDomainEvent) GetID() uuid.UUID {
	return b.id
}

func (b *baseOrderDomainEvent) GetName() string {
	return b.name
}

func (b *baseOrderDomainEvent) OccurredAt() time.Time {
	return b.occurredAt
}

func (b *baseOrderDomainEvent) OrderId() uuid.UUID {
	return b.orderId
}
