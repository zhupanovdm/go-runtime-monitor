package storage

import (
	"context"

	"github.com/zhupanovdm/go-runtime-monitor/config"
	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
)

type (
	// Storage is representing metrics operations on storage.
	Storage interface {
		// IsPersistent returns true if the storage is persistent, otherwise false.
		IsPersistent() bool

		// Init initializes storage. Should be called before working with storage.
		Init(ctx context.Context) error

		// Ping checks if storage is online.
		Ping(ctx context.Context) error

		// Close releases storage resources. Should be called on end of working with storage.
		Close(ctx context.Context)

		// Get queries single metric. Returns nil if metric not found.
		Get(ctx context.Context, id string, typ metric.Type) (*metric.Metric, error)

		// GetAll queries all metrics.
		GetAll(ctx context.Context) (metric.List, error)

		// Update registers or updates single metric.
		Update(ctx context.Context, mtr *metric.Metric) error

		// UpdateBulk registers or updates all metrics in list.
		UpdateBulk(ctx context.Context, list metric.List) error

		// Clear deletes all metrics in storage.
		Clear(ctx context.Context) error
	}

	// Factory produces initialized storage object.
	Factory func(*config.Config) Storage
)

// New returns new metrics storage created with first factory that will return non-nil value.
func New(cfg *config.Config, factories ...Factory) Storage {
	for _, create := range factories {
		if storage := create(cfg); storage != nil {
			return storage
		}
	}
	return nil
}
