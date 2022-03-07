package agent

import (
	"context"

	"github.com/stretchr/testify/require"

	"github.com/zhupanovdm/go-runtime-monitor/model/metric"
	"github.com/zhupanovdm/go-runtime-monitor/pkg/task"
)

var _ ReporterService = (*stubReporter)(nil)

type stubReporter struct {
	t        require.TestingT
	consumer func(*metric.Metric)
}

func (s *stubReporter) Report(context.Context) error {
	return nil
}

func (s *stubReporter) ReportBulk(context.Context) error {
	return nil
}

func (s *stubReporter) Name() string {
	return "Stub Reporter"
}

func (s *stubReporter) BackgroundTask() task.Task {
	return task.VoidTask
}

func (s *stubReporter) Publish(_ context.Context, mtr *metric.Metric) {
	v, err := mtr.Type().New()
	require.NoErrorf(s.t, err, "error creating zero value of type: %v", mtr.Type())

	s.consumer(&metric.Metric{
		ID:    mtr.ID,
		Value: v,
	})
}

func NewStubReporter(t require.TestingT, consumer func(*metric.Metric)) ReporterService {
	return &stubReporter{
		t:        t,
		consumer: consumer,
	}
}
