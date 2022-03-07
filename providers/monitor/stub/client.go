package stub

import (
	"context"

	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/providers/monitor"
)

var _ monitor.Provider = (*client)(nil)

type client struct{}

func (c *client) Update(context.Context, *metric.Metric) error {
	return nil
}

func (c *client) UpdateBulk(context.Context, metric.List) error {
	return nil
}

func (c *client) Value(_ context.Context, _ string, typ metric.Type) (metric.Value, error) {
	return typ.New()
}

func New() monitor.Provider {
	return &client{}
}
