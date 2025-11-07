package jobs

import (
	"context"
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/pkg/errs"
	"github.com/labstack/gommon/log"
	"github.com/robfig/cron/v3"
)

var _ cron.Job = &moveCouriersJob{}

type moveCouriersJob struct {
	handler commands.MoveCouriersCommandHandler
}

func NewMoveCouriersJob(handler commands.MoveCouriersCommandHandler) (cron.Job, error) {
	if handler == nil {
		return nil, errs.NewValueIsRequiredError("handler")
	}

	return &moveCouriersJob{handler: handler}, nil
}

func (job *moveCouriersJob) Run() {
	err := job.handler.Handle(context.Background())
	if err != nil {
		log.Error(err)
	}
}
