package agent

import (
	"context"

	"github.com/zhupanovdm/go-runtime-monitor/pkg"
)

type (
	// CollectorService is an application service responsible for gathering and publishing runtime metrics collection.
	CollectorService interface {
		pkg.BackgroundService

		// Poll collects actual runtime metrics
		Poll(context.Context)
	}

	// Collector is metrics collecting strategy.
	Collector func(ctx context.Context, froze *Froze) error

	// ReporterService is an application service responsible for transporting published metrics on external monitor service.
	ReporterService interface {
		pkg.BackgroundService

		// Report transmits metrics to server
		Report(context.Context)
	}
)
