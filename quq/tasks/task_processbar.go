package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

const (
	TypeProcessBar = "tasks:processbar"
)

type ProcessBarPayload = struct {
	Count int64
}

func NewProcessBarTask(count int64) (*asynq.Task, error) {
	payload, err := json.Marshal(ProcessBarPayload{Count: count})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeProcessBar, payload), nil
}

func HandleProcessBarTask(ctx context.Context, t *asynq.Task) error {
	processStep := make(chan int64, 100)
	errorQueue := make(chan interface{})
	cancelQueue := make(chan bool)
	go handleProcessBarTask(ctx, t, errorQueue, processStep, cancelQueue)
	var err interface{}
	for {
		select {
		case <-ctx.Done():
			Logger.Info("canceled")
			cancelQueue <- true
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

func handleProcessBarTask(
	ctx context.Context,
	t *asynq.Task,
	errorQueue chan interface{},
	processStep chan int64,
	cancelQueue chan bool) {
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
		select {
		case <-cancelQueue:
			Logger.Info("canceled recived")
			return
		default:
		}
		processStep <- p.Count
	}

}
