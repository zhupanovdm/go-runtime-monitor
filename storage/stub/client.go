package stub

import (
	"context"

	"github.com/zhupanovdm/go-runtime-monitor/config"
	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/storage"
)

var _ storage.Storage = (*client)(nil)

type client struct{}

func (c *client) Clear(context.Context) error {
	return nil
}

func (c *client) IsPersistent() bool {
	return false
}

func (c *client) Init(context.Context) error {
	return nil
}

func (c *client) Ping(context.Context) error {
	return nil
}

func (c *client) Close(context.Context) {}

func (c *client) Get(context.Context, string, metric.Type) (*metric.Metric, error) {
	return nil, nil
}

func (c *client) GetAll(context.Context) (metric.List, error) {
	return make(metric.List, 0), nil
}

func (c *client) Update(context.Context, string, metric.Value) error {
	return nil
}

func (c *client) UpdateBulk(context.Context, metric.List) error {
	return nil
}

func New(*config.Config) storage.Storage {
	return &client{}
}
