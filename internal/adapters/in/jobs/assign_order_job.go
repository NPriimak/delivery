package jobs

import (
	"context"
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/pkg/errs"
	"github.com/labstack/gommon/log"
	"github.com/robfig/cron/v3"
)

var _ cron.Job = &AssignOrderJob{}

type AssignOrderJob struct {
	assignOrdersCommandHandler commands.AssignOrderCommandHandler
}

func NewAssignOrderJob(
	assignOrdersCommandHandler commands.AssignOrderCommandHandler) (*AssignOrderJob, error) {
	if assignOrdersCommandHandler == nil {
		return nil, errs.NewValueIsRequiredError("moveCouriersCommandHandler")
	}

	return &AssignOrderJob{
		assignOrdersCommandHandler: assignOrdersCommandHandler}, nil
}

func (j *AssignOrderJob) Run() {
	log.Info("Assign order")
	ctx := context.Background()
	command := commands.NewAssignOrdersCommand()
	err := j.assignOrdersCommandHandler.Handle(ctx, command)
	if err != nil {
		log.Error(err)
	}
}
