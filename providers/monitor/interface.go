package monitor

import (
	"context"

	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
)

type Provider interface {
	Update(ctx context.Context, mtr *metric.Metric) error
	Value(ctx context.Context, id string, typ metric.Type) (metric.Value, error)
}
