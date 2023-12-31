package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"

	"quq/tasks"

	"github.com/hibiken/asynq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sys/unix"
)

const redisAddr = "127.0.0.1:36379"

func main() {
	// metric
	httpServerMux := http.NewServeMux()
	httpServerMux.Handle("/metrics", promhttp.Handler())
	metricsSrv := &http.Server{
		Addr:    ":7070",
		Handler: httpServerMux,
	}
	done := make(chan struct{})

	go func() {
		err := metricsSrv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Printf("Error: metrics server error: %v", err)
		}
		close(done)
	}()

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			ErrorHandler: tasks.ErrorHandler{},
			Logger:       tasks.Logger,
			LogLevel:     asynq.DebugLevel,
			IsFailure: func(err error) bool {
				if errors.Is(err, context.Canceled) {
					tasks.Logger.Info("IsFailure: ", err, false)
					return false
				}
				tasks.Logger.Info("IsFailure: ", err, true)
				return true
			},
		},
	)

	mux := asynq.NewServeMux()
	mux.Use(metricsMiddleware)
	mux.HandleFunc(tasks.TypeProcessBar, tasks.HandleProcessBarTask)

	if err := srv.Start(mux); err != nil {
		log.Fatalf("could not start server: %v", err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, unix.SIGTERM, unix.SIGINT)
	<-sigs

	srv.Shutdown()

	if err := metricsSrv.Shutdown(context.Background()); err != nil {
		log.Printf("Error: metrics server shutdown error: %v", err)
	}
	<-done
}
