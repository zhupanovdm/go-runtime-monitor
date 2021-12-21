package service

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestCounter(t *testing.T) {
	assert.Equal(t, "counter/foo/10", counter("foo", 10))
}

func TestGauge(t *testing.T) {
	assert.Equal(t, "gauge/foo/0.100000", gauge("foo", 0.1))
}

func TestGaugeu(t *testing.T) {
	assert.Equal(t, "gauge/foo/1.000000", gaugeu("foo", 1))
}

func TestPollRuntimeMetrics(t *testing.T) {
	data := make(chan string, 128)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		pollRuntimeMetrics(data, 0)
		wg.Done()
	}()

	wg.Wait()
	assert.Equal(t, 27, len(data))
}
