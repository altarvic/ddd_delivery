package order

import (
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"
	"github.com/google/uuid"
	"reflect"
	"time"
)

var _ ddd.DomainEvent = &CompletedDomainEvent{}

type CompletedDomainEvent struct {
	Id         uuid.UUID
	Name       string
	OccurredAt time.Time

	OrderId   uuid.UUID
	CourierId uuid.UUID

	isValid bool
}

func NewCompletedDomainEvent(orderId uuid.UUID, courierId uuid.UUID) (ddd.DomainEvent, error) {

	event := &CompletedDomainEvent{}
	if orderId == uuid.Nil {
		return event, errs.NewValueIsRequiredError("orderId")
	}

	if courierId == uuid.Nil {
		return event, errs.NewValueIsRequiredError("courierId")
	}

	event.Id = uuid.New()
	event.Name = reflect.TypeOf(event).Elem().Name()
	event.OccurredAt = time.Now().UTC()
	event.OrderId = orderId
	event.CourierId = courierId
	event.isValid = true

	return event, nil
}

func (e *CompletedDomainEvent) GetID() uuid.UUID {
	return e.Id
}

func (e *CompletedDomainEvent) GetName() string {
	return e.Name
}

func (e *CompletedDomainEvent) IsValid() bool {
	return e.isValid
}
