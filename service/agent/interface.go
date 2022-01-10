package agent

import (
	"context"

	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/pkg"
)

type CollectorService interface {
	pkg.BackgroundService
}

type ReporterService interface {
	pkg.BackgroundService
	Publish(ctx context.Context, mtr *metric.Metric)
}
