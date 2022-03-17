package monitor

import (
	"context"

	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/pkg"
)

// Monitor application service is responsible for operations with metrics.
type Monitor interface {
	pkg.BackgroundService

	// Restore restores previously dumped metrics.
	Restore(ctx context.Context) error

	// Get queries single metric. Will return nil if requested metric not found.
	Get(ctx context.Context, id string, typ metric.Type) (*metric.Metric, error)

	// GetAll queries all registered metrics.
	GetAll(ctx context.Context) (metric.List, error)

	// Update registers or updates previously registered metric.
	Update(ctx context.Context, mtr *metric.Metric) error

	// UpdateBulk registers or updates previously registered metrics in list.
	UpdateBulk(ctx context.Context, list metric.List) error

	// Ping diagnoses service state.
	Ping(ctx context.Context) error
}
