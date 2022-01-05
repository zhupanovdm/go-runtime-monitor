package monitor

import (
	"context"

	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/pkg"
)

type Service interface {
	pkg.Service

	Get(ctx context.Context, id string, typ metric.Type) (*metric.Metric, error)
	GetAll(ctx context.Context) ([]*metric.Metric, error)

	Update(ctx context.Context, mtr *metric.Metric) error
}
