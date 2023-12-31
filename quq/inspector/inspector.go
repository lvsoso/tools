package main

import (
	"flag"
	"fmt"

	"github.com/hibiken/asynq"
)

var (
	// create bool
	cancel bool
	delete bool
	get    bool
	// retry  bool

	queue     string
	redisAddr string
	taskId    string
)

func init() {
	// flag.BoolVar(&create, "create", false, "create")
	flag.BoolVar(&get, "g", false, "get")
	flag.BoolVar(&cancel, "c", false, "cancel")
	flag.BoolVar(&delete, "d", false, "delete")
	// flag.BoolVar(&retry, "retry", false, "retry")
	flag.StringVar(&queue, "q", "default", "queue")
	flag.StringVar(&redisAddr, "r", "127.0.0.1:36379", "redis uri")
	flag.StringVar(&taskId, "t", "", "task id")
}

func main() {
	flag.Parse()

	inspector := asynq.NewInspector(asynq.RedisClientOpt{Addr: redisAddr})
	defer inspector.Close()

	if taskId == "" {
		println("empty task id")
		return
	}

	if get {
		tinfo, err := inspector.GetTaskInfo(queue, taskId)
		if err != nil {
			fmt.Printf("get info error : %+v\n", err)
			return
		}
		fmt.Printf("Tindo: %++v\n", tinfo)
	}

	if cancel {
		// CancelProcessing is best-effort
		err := inspector.CancelProcessing(taskId)
		if err != nil {
			fmt.Printf("cancel error : %+v\n", err)
			return
		}
	}

	if delete {
		err := inspector.DeleteTask(queue, taskId)
		if err != nil {
			fmt.Printf("delete error : %+v\n", err)
			return
		}
	}
}
