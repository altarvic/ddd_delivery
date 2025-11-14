package eventhandlers

import (
	"context"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"
)

var _ ddd.EventHandler = &orderEventsHandler{}

type orderEventsHandler struct {
	notificationProducer ports.NotificationProducer
}

func NewOrderDomainEventsHandler(notificationProducer ports.NotificationProducer) (ddd.EventHandler, error) {
	if notificationProducer == nil {
		return nil, errs.NewValueIsRequiredError("notificationProducer")
	}

	return &orderEventsHandler{
		notificationProducer: notificationProducer,
	}, nil
}

func (eh *orderEventsHandler) Handle(ctx context.Context, domainEvent ddd.DomainEvent) error {
	return eh.notificationProducer.Publish(ctx, domainEvent)
}
