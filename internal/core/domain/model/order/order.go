package order

import (
	"delivery/internal/core/domain/kernel"
	"delivery/internal/pkg/ddd"
	"errors"
	"github.com/google/uuid"
)

type Order struct {
	id        uuid.UUID
	courierId *uuid.UUID
	location  kernel.Location
	volume    int
	status    Status

	events []ddd.DomainEvent
}

func NewOrder(orderId uuid.UUID, location kernel.Location, volume int) (*Order, error) {
	if orderId == uuid.Nil {
		return nil, errors.New("empty orderId")
	}

	if location.IsEmpty() {
		return nil, errors.New("empty location")
	}

	if volume <= 0 {
		return nil, errors.New("volume <= 0")
	}

	orderCreatedEvent, err := NewCreatedDomainEvent(orderId)
	if err != nil {
		return nil, err
	}

	order := &Order{
		id:       orderId,
		location: location,
		volume:   volume,
		status:   StatusCreated,
		events:   []ddd.DomainEvent{},
	}

	order.RaiseDomainEvent(orderCreatedEvent)

	return order, nil
}

func (o *Order) GetDomainEvents() []ddd.DomainEvent {
	return o.events
}

func (o *Order) ClearDomainEvents() {
	o.events = []ddd.DomainEvent{}
}

func (o *Order) RaiseDomainEvent(event ddd.DomainEvent) {
	o.events = append(o.events, event)
}

func (o *Order) Id() uuid.UUID {
	return o.id
}

func (o *Order) CourierId() *uuid.UUID {
	return o.courierId
}

func (o *Order) Location() kernel.Location {
	return o.location
}

func (o *Order) Volume() int {
	return o.volume
}

func (o *Order) Status() Status {
	return o.status
}

func (o *Order) Equals(other *Order) bool {
	return other != nil && o.id == other.id
}

func (o *Order) AssignCourier(courierId uuid.UUID) error {
	if o.status != StatusCreated {
		return errors.New("courier already assigned")
	}

	if courierId == uuid.Nil {
		return errors.New("empty courierId")
	}

	o.courierId = &courierId
	o.status = StatusAssigned
	return nil
}

func (o *Order) Complete() error {
	if o.status == StatusCompleted {
		return errors.New("already completed")
	}

	if o.status == StatusCreated {
		return errors.New("order w/o assigned courier")
	}

	o.status = StatusCompleted

	orderCompletedEvent, err := NewCompletedDomainEvent(o.id, *o.courierId)
	if err != nil {
		return err
	}

	o.RaiseDomainEvent(orderCompletedEvent)

	return nil
}

// RestoreOrder should be used ONLY inside Repository
func RestoreOrder(id uuid.UUID, courierId *uuid.UUID, location kernel.Location, volume int, status Status) *Order {
	return &Order{
		id:        id,
		courierId: courierId,
		location:  location,
		volume:    volume,
		status:    status,
	}
}
