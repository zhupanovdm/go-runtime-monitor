package storage

import (
	"context"

	"github.com/zhupanovdm/go-runtime-monitor/config"
	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
)

type Provider func(*config.Config) Storage

func New(cfg *config.Config, providers ...Provider) Storage {
	for _, provider := range providers {
		if s := provider(cfg); s != nil {
			return s
		}
	}
	return nil
}

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
	Init(ctx context.Context) error
	Ping(ctx context.Context) error
	Close(ctx context.Context)

	GetAll(ctx context.Context) (metric.List, error)
	UpdateBulk(ctx context.Context, list metric.List) error
}
