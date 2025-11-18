package order

import (
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"
	"github.com/google/uuid"
	"reflect"
	"time"
)

var _ ddd.DomainEvent = &CreatedDomainEvent{}

type CreatedDomainEvent struct {
	Id         uuid.UUID
	Name       string
	OccurredAt time.Time

	OrderId uuid.UUID

	isValid bool
}

func NewCreatedDomainEvent(orderId uuid.UUID) (ddd.DomainEvent, error) {
	event := &CreatedDomainEvent{}
	if orderId == uuid.Nil {
		return event, errs.NewValueIsRequiredError("orderId")
	}

	event.Id = uuid.New()
	event.Name = reflect.TypeOf(event).Elem().Name()
	event.OccurredAt = time.Now().UTC()
	event.OrderId = orderId
	event.isValid = true

	return event, nil
}

func (e *CreatedDomainEvent) GetID() uuid.UUID {
	return e.Id
}

func (e *CreatedDomainEvent) GetName() string {
	return e.Name
}

func (e *CreatedDomainEvent) IsValid() bool {
	return e.isValid
}
