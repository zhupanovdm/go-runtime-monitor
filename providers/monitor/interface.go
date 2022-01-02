package monitor

import (
	"context"
	metric2 "github.com/zhupanovdm/go-runtime-monitor/model/metric"
)

type Provider interface {
	Update(ctx context.Context, mtr *metric2.Metric) error
	Value(ctx context.Context, id string, typ metric2.Type) (metric2.Value, error)
}
