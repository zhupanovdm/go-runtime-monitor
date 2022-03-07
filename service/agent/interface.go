package agent

import (
	"context"

	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/pkg"
)

type CollectorService interface {
	pkg.BackgroundService
	Poll(context.Context)
}

type ReporterService interface {
	pkg.BackgroundService
	Publish(context.Context, *metric.Metric)

	Report(context.Context) error
	ReportBulk(context.Context) error
}
