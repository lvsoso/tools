package main

import (
	"log"

	"quq/tasks"

	"github.com/hibiken/asynq"
)

const redisAddr = "127.0.0.1:36379"

func main() {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer client.Close()

	task, err := tasks.NewProcessBarTask(30)
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}
	info, err := client.Enqueue(task, asynq.MaxRetry(0))
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

}
