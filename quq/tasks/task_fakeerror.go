package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

const (
	TypeFakeError = "tasks:fakeerror"
)

type FakeErrorPayload = struct {
	Count int64
}

func NewFakeErrorTask(count int64) (*asynq.Task, error) {
	payload, err := json.Marshal(ProcessBarPayload{Count: count})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeProcessBar, payload), nil
}

func HandleFakeErrorTask(ctx context.Context, t *asynq.Task) error {
	processStep := make(chan int64, 100)
	errorQueue := make(chan interface{})
	go handleFakeErrorTask(ctx, t, errorQueue, processStep)
	var err interface{}
	for {
		select {
		case <-ctx.Done():
			Logger.Info("canceled")
			return nil
		case err = <-errorQueue:
			if err != nil {
				Logger.Error(err)
				return fmt.Errorf("unknown %v", err)
			}
			Logger.Info("finnished")
			return nil
		case step := <-processStep:
			Logger.Info(fmt.Sprintf("step %d \n", step))
		}
	}
}

func handleFakeErrorTask(ctx context.Context, t *asynq.Task, errorQueue chan interface{}, processStep chan int64) {
	defer func() {
		if err := recover(); err != nil {
			Logger.Error(err)
			errorQueue <- err
		}
	}()
	var p ProcessBarPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		errorQueue <- fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
		return
	}

	for i := 0; i < int(p.Count); i++ {
		time.Sleep(1 * time.Second)
		Logger.Info(fmt.Sprintf("Count=%d", i))

		if i == 5 {
			panic(errors.New("abc"))
		}

		processStep <- p.Count
	}

}
