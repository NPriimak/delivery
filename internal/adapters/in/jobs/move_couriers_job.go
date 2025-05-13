package jobs

import (
	"context"
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/pkg/errs"
	"github.com/labstack/gommon/log"
	"github.com/robfig/cron/v3"
)

var _ cron.Job = &MovingCouriersJob{}

type MovingCouriersJob struct {
	moveCourierCommandHandler commands.MoveCouriersCommandHandler
}

func NewMoveCouriersJob(commandHandler commands.MoveCouriersCommandHandler) (*MovingCouriersJob, error) {
	if commandHandler == nil {
		return nil, errs.NewValueIsRequiredError("moveCouriersCommandHandler")
	}

	return &MovingCouriersJob{
		moveCourierCommandHandler: commandHandler}, nil
}

func (j *MovingCouriersJob) Run() {
	log.Info("Moving couriers")
	ctx := context.Background()
	command, err := commands.NewMoveCouriersCmd()
	if err != nil {
		log.Error(err)
	}
	err = j.moveCourierCommandHandler.Handle(ctx, command)
	if err != nil {
		log.Error(err)
	}
}
