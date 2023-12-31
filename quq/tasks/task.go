package tasks

import (
	"context"

	"github.com/hibiken/asynq"
)

type ErrorHandler struct {
}

func (eh ErrorHandler) HandleError(ctx context.Context, task *asynq.Task, err error) {

	Logger.Info("HandleError ", err.Error())
}
