package kafka

import (
	"context"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/core/ports"
	"delivery/internal/generated/queues/orderpb"
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ ports.NotificationProducer = &orderChangedNotificationProducer{}

type orderChangedNotificationProducer struct {
	topic    string
	producer sarama.SyncProducer
}

func NewOrderChangedNotificationProducer(brokers []string, topic string) (ports.NotificationProducer, error) {
	if len(brokers) == 0 {
		return nil, errs.NewValueIsRequiredError("brokers")
	}
	if topic == "" {
		return nil, errs.NewValueIsRequiredError("topic")
	}

	saramaCfg := sarama.NewConfig()
	saramaCfg.Version = sarama.V3_9_1_0
	saramaCfg.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, saramaCfg)
	if err != nil {
		return nil, fmt.Errorf("create sync producer: %w", err)
	}

	return &orderChangedNotificationProducer{
		topic:    topic,
		producer: producer,
	}, nil
}

func (p *orderChangedNotificationProducer) Close() error {
	return p.producer.Close()
}

func (p *orderChangedNotificationProducer) Publish(ctx context.Context, domainEvent ddd.DomainEvent) error {

	var (
		integrationEvent any
		key              string
	)

	switch domainEvent.(type) {
	case *order.CreatedDomainEvent:
		createdEvent := domainEvent.(*order.CreatedDomainEvent)
		integrationEvent = p.mapCreatedDomainEventToIntegrationEvent(createdEvent)
		key = createdEvent.OrderId.String()
	case *order.CompletedDomainEvent:
		completedEvent := domainEvent.(*order.CompletedDomainEvent)
		integrationEvent = p.mapCompletedDomainEventToIntegrationEvent(completedEvent)
		key = completedEvent.OrderId.String()
	default:
		return errors.New("unknown order changed event type")
	}

	data, err := json.Marshal(integrationEvent)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(data),
	}

	resultCh := make(chan error, 1)

	go func() {
		_, _, err := p.producer.SendMessage(msg)
		resultCh <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-resultCh:
		return err
	}
}

func (p *orderChangedNotificationProducer) mapCreatedDomainEventToIntegrationEvent(domainEvent *order.CreatedDomainEvent) *orderpb.OrderCreatedIntegrationEvent {
	return &orderpb.OrderCreatedIntegrationEvent{
		EventId:    domainEvent.GetID().String(),
		EventType:  domainEvent.GetName(),
		OccurredAt: timestamppb.New(domainEvent.OccurredAt),
		OrderId:    domainEvent.OrderId.String(),
	}
}

func (p *orderChangedNotificationProducer) mapCompletedDomainEventToIntegrationEvent(domainEvent *order.CompletedDomainEvent) *orderpb.OrderCompletedIntegrationEvent {
	return &orderpb.OrderCompletedIntegrationEvent{
		EventId:    domainEvent.GetID().String(),
		EventType:  domainEvent.GetName(),
		OccurredAt: timestamppb.New(domainEvent.OccurredAt),
		OrderId:    domainEvent.OrderId.String(),
		CourierId:  domainEvent.CourierId.String(),
	}
}
