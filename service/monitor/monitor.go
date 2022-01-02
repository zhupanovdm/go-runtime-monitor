package monitor

import (
	"context"
	"fmt"

	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/storage"
)

var _ MetricsMonitorService = (*monitorSvc)(nil)

type monitorSvc struct {
	gauge   storage.GaugeStorage
	counter storage.CounterStorage
}

func (m *monitorSvc) Save(ctx context.Context, mtr *metric.Metric) (err error) {
	switch value := mtr.Value.(type) {
	case *metric.Gauge:
		err = m.gauge.Update(ctx, mtr.ID, *value)
	case *metric.Counter:
		err = m.counter.Update(ctx, mtr.ID, *value)
	default:
		err = fmt.Errorf("type is not supported yet: %T", value)
	}
	return
}

func (m *monitorSvc) Get(ctx context.Context, id string, typ metric.Type) (value *metric.Metric, err error) {
	switch typ {
	case metric.GaugeType:
		return m.gauge.Get(ctx, id)
	case metric.CounterType:
		return m.counter.Get(ctx, id)
	default:
		err = fmt.Errorf("unknown metric type: %v", typ)
	}
	return
}

func (m *monitorSvc) GetAll(ctx context.Context) (list []*metric.Metric, err error) {
	if list, err = m.gauge.GetAll(ctx); err != nil {
		return
	}
	var counters []*metric.Metric
	if counters, err = m.counter.GetAll(ctx); err != nil {
		return
	}

	list = append(list, counters...)
	return
}

func NewMetricsMonitor(gaugeStorage storage.GaugeStorage, counterStorage storage.CounterStorage) MetricsMonitorService {
	return &monitorSvc{
		gauge:   gaugeStorage,
		counter: counterStorage,
	}
}
