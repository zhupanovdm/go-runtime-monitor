package agent

import (
	"context"
	"github.com/zhupanovdm/go-runtime-monitor/model/metric"

	"github.com/zhupanovdm/go-runtime-monitor/pkg/task"
)

type CollectorService interface {
	BackgroundRunner
}

type ReporterService interface {
	BackgroundRunner
	Publish(ctx context.Context, mtr *metric.Metric)
}

type BackgroundRunner interface {
	BackgroundTask() task.Task
}
