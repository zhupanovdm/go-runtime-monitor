package storage

import (
	"context"

	"github.com/zhupanovdm/go-runtime-monitor/config"
	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
)

type Factory func(*config.Config) Storage

func New(cfg *config.Config, factories ...Factory) Storage {
	for _, create := range factories {
		if storage := create(cfg); storage != nil {
			return storage
		}
	}
	return nil
}

type Storage interface {
	IsPersistent() bool

	Init(ctx context.Context) error
	Ping(ctx context.Context) error
	Close(ctx context.Context)

	Get(ctx context.Context, id string, typ metric.Type) (*metric.Metric, error)
	GetAll(ctx context.Context) (metric.List, error)

	Update(ctx context.Context, id string, value metric.Value) error
	UpdateBulk(ctx context.Context, list metric.List) error

	Clear(ctx context.Context) error
}
