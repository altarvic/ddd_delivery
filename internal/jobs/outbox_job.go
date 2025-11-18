package jobs

import (
	"context"
	outb "delivery/internal/adapters/out/outbox"
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"
	"delivery/internal/pkg/outbox"
	"github.com/labstack/gommon/log"
	"github.com/robfig/cron/v3"
	"time"
)

var _ cron.Job = &moveCouriersJob{}

type outboxJob struct {
	ob       outb.OutboxRepository
	mediatr  ddd.Mediatr
	registry outbox.EventRegistry
}

func NewOutboxJob(ob outb.OutboxRepository, mediatr ddd.Mediatr, registry outbox.EventRegistry) (cron.Job, error) {
	if ob == nil {
		return nil, errs.NewValueIsRequiredError("ob")
	}

	if mediatr == nil {
		return nil, errs.NewValueIsRequiredError("mediatr")
	}

	if registry == nil {
		return nil, errs.NewValueIsRequiredError("registry")
	}

	return &outboxJob{
		ob:       ob,
		mediatr:  mediatr,
		registry: registry,
	}, nil
}

func (job *outboxJob) Run() {

	var messages []*outbox.Message
	var err error
	ctx := context.Background()

	messages, err = job.ob.GetNotPublishedMessages(ctx)
	if err != nil {
		log.Error(err)
		return
	}

	if messages == nil || len(messages) == 0 {
		return
	}

	processedMessages := make([]*outbox.Message, 0, len(messages))
	for _, msg := range messages {

		domainEvent, err := job.registry.DecodeDomainEvent(msg)
		if err != nil {
			log.Error(err)
			continue
		}
		log.Info(domainEvent)

		err = job.mediatr.Publish(ctx, domainEvent)
		if err != nil {
			log.Error(err)
			continue
		}

		now := time.Now().UTC()
		msg.ProcessedAtUtc = &now
		processedMessages = append(processedMessages, msg)
	}

	if len(processedMessages) > 0 {
		err = job.ob.Save(ctx, processedMessages...)
		if err != nil {
			log.Error(err)
		}
	}
}
