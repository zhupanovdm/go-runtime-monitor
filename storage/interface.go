package storage

import (
	"context"

	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
)

type GaugeStorage interface {
	Get(ctx context.Context, id string) (*metric.Metric, error)
	GetAll(ctx context.Context) ([]*metric.Metric, error)

	Update(ctx context.Context, id string, value metric.Gauge) error
}

type CounterStorage interface {
	Get(ctx context.Context, id string) (*metric.Metric, error)
	GetAll(ctx context.Context) ([]*metric.Metric, error)

	Update(ctx context.Context, id string, value metric.Counter) error
}
