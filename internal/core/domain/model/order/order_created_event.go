package order

import (
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"
	"github.com/google/uuid"
	"reflect"
)

var _ ddd.DomainEvent = &CreatedDomainEvent{}

type CreatedDomainEvent struct {
	baseOrderDomainEvent

	isValid bool
}

func NewCreatedDomainEvent(orderId uuid.UUID) (ddd.DomainEvent, error) {
	event := &CreatedDomainEvent{}
	if orderId == uuid.Nil {
		return event, errs.NewValueIsRequiredError("orderId")
	}

	event.baseOrderDomainEvent = newBaseOrderDomainEvent(reflect.TypeOf(event).Elem().Name(), orderId)
	event.isValid = true

	return event, nil
}

func (o *CreatedDomainEvent) IsValid() bool {
	return o.isValid
}
