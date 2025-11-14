package order

import (
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"
	"github.com/google/uuid"
	"reflect"
)

var _ ddd.DomainEvent = &CompletedDomainEvent{}

type CompletedDomainEvent struct {
	baseOrderDomainEvent

	courierId uuid.UUID
	isValid   bool
}

func NewCompletedDomainEvent(orderId uuid.UUID, courierId uuid.UUID) (ddd.DomainEvent, error) {

	event := &CompletedDomainEvent{}
	if orderId == uuid.Nil {
		return event, errs.NewValueIsRequiredError("orderId")
	}

	if courierId == uuid.Nil {
		return event, errs.NewValueIsRequiredError("courierId")
	}

	event.baseOrderDomainEvent = newBaseOrderDomainEvent(reflect.TypeOf(event).Elem().Name(), orderId)
	event.courierId = courierId
	event.isValid = true

	return event, nil
}

func (o *CompletedDomainEvent) CourierId() uuid.UUID {
	return o.courierId
}

func (o *CompletedDomainEvent) IsValid() bool {
	return o.isValid
}
