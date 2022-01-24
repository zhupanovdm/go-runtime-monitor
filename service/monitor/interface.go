package monitor

import (
	"context"

	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/pkg"
)

type Monitor interface {
	pkg.BackgroundService

	Restore(ctx context.Context) error

	Get(ctx context.Context, id string, typ metric.Type) (*metric.Metric, error)
	GetAll(ctx context.Context) (metric.List, error)

	Update(ctx context.Context, mtr *metric.Metric) error

	Ping(ctx context.Context) error
}
