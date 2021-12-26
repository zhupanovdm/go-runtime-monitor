package service

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zhupanovdm/go-runtime-monitor/internal/model/metric"
)

func TestPollRuntimeMetrics(t *testing.T) {
	data := make(chan *metric.Metric, 128)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		pollRuntimeMetrics(data, 0)
		wg.Done()
	}()

	wg.Wait()
	assert.Equal(t, 27, len(data))
}
