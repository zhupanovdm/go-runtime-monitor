package monitor

import (
	"context"

	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
)

type MetricsMonitorService interface {
	Get(ctx context.Context, id string, typ metric.Type) (*metric.Metric, error)
	GetAll(ctx context.Context) ([]*metric.Metric, error)

	Save(ctx context.Context, mtr *metric.Metric) error
}
