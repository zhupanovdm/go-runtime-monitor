package agent

import (
	"context"

	"github.com/zhupanovdm/go-runtime-monitor/pkg"
)

type CollectorService interface {
	pkg.BackgroundService
	Poll(context.Context)
}

type ReporterService interface {
	pkg.BackgroundService
	Report(context.Context)
}
