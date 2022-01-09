package storage

import (
	"context"

	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
)

type GaugeStorage interface {
	Storage

	Get(ctx context.Context, id string) (*metric.Metric, error)
	Update(ctx context.Context, id string, value metric.Gauge) error
}

type CounterStorage interface {
	Storage

	Get(ctx context.Context, id string) (*metric.Metric, error)
	Update(ctx context.Context, id string, value metric.Counter) error
}

type Storage interface {
	GetAll(ctx context.Context) (metric.List, error)
	UpdateBulk(ctx context.Context, list metric.List) error
}
