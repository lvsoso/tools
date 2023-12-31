package main

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics variables.
var (
	processedCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "processed_tasks_total",
			Help: "The total number of processed tasks",
		},
		[]string{"task_type"},
	)

	failedCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "failed_tasks_total",
			Help: "The total number of times processing failed",
		},
		[]string{"task_type"},
	)

	inProgressGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "in_progress_tasks",
			Help: "The number of tasks currently being processed",
		},
		[]string{"task_type"},
	)
)

func metricsMiddleware(next asynq.Handler) asynq.Handler {
	return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
		inProgressGauge.WithLabelValues(t.Type()).Inc()
		err := next.ProcessTask(ctx, t)
		inProgressGauge.WithLabelValues(t.Type()).Dec()
		if err != nil {
			failedCounter.WithLabelValues(t.Type()).Inc()
		}
		processedCounter.WithLabelValues(t.Type()).Inc()
		return err
	})
}
