package jobs

import (
	"context"
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/pkg/errs"
	"github.com/labstack/gommon/log"
	"github.com/robfig/cron/v3"
)

var _ cron.Job = &assignOrdersJob{}

type assignOrdersJob struct {
	handler commands.AssignOrderCommandHandler
}

func NewAssignOrdersJob(handler commands.AssignOrderCommandHandler) (cron.Job, error) {
	if handler == nil {
		return nil, errs.NewValueIsRequiredError("handler")
	}

	return &assignOrdersJob{handler: handler}, nil
}

func (job *assignOrdersJob) Run() {
	err := job.handler.Handle(context.Background())
	if err != nil {
		log.Error(err)
	}
}
