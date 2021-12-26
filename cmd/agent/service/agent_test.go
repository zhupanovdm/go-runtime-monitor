package service

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zhupanovdm/go-runtime-monitor/internal/model/metric"
)

var ctx = context.TODO()
var voidFunc = func() {}

var foo = &metric.Metric{Id: "foo"}
var bar = &metric.Metric{Id: "bar"}
var baz = &metric.Metric{Id: "baz"}

func TestPublish(t *testing.T) {
	data := make(chan *metric.Metric)
	e := publish(data, func(c chan<- *metric.Metric) {
		c <- foo
		c <- bar
		c <- baz
	})

	go func() {
		e.Start()
		e.Exec(ctx, voidFunc)
		e.End()
	}()

	assert.Equal(t, <-data, foo)
	assert.Equal(t, <-data, bar)
	assert.Equal(t, <-data, baz)

	_, ok := <-data
	assert.False(t, ok)
}

func TestSubscribe(t *testing.T) {
	data := make(chan *metric.Metric, 3)
	data <- foo
	data <- bar
	data <- baz

	coll := make([]*metric.Metric, 0, 3)
	e := subscribe(data, func(m *metric.Metric) error {
		coll = append(coll, m)
		return nil
	})

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		e.Start()
		e.Exec(ctx, voidFunc)
		e.End()
		wg.Done()
	}()

	wg.Wait()
	close(data)

	assert.ElementsMatch(t, coll, []*metric.Metric{foo, bar, baz})
}

func TestSubscribePartial(t *testing.T) {
	data := make(chan *metric.Metric, 3)
	data <- foo
	data <- bar

	coll := make([]*metric.Metric, 0, 3)
	e := subscribe(data, func(m *metric.Metric) error {
		coll = append(coll, m)
		return nil
	})

	var wg sync.WaitGroup
	poller := func() {
		e.Start()
		e.Exec(ctx, voidFunc)
		e.End()
		wg.Done()
	}

	wg.Add(1)
	go poller()
	wg.Wait()
	assert.ElementsMatch(t, coll, []*metric.Metric{foo, bar})

	data <- baz
	wg.Add(1)
	go poller()
	wg.Wait()
	close(data)

	assert.ElementsMatch(t, coll, []*metric.Metric{foo, bar, baz})
}

func TestSubscribeClosed(t *testing.T) {
	data := make(chan *metric.Metric, 3)
	data <- foo
	data <- bar

	coll := make([]*metric.Metric, 0, 3)
	e := subscribe(data, func(s *metric.Metric) error {
		coll = append(coll, s)
		return nil
	})

	var wg sync.WaitGroup
	poller := func() {
		e.Start()
		e.Exec(ctx, voidFunc)
		e.End()
		wg.Done()
	}

	wg.Add(1)
	go poller()
	wg.Wait()
	assert.ElementsMatch(t, coll, []*metric.Metric{foo, bar})

	close(data)
	wg.Add(1)
	go poller()
	wg.Wait()

	assert.ElementsMatch(t, coll, []*metric.Metric{foo, bar})
}

func TestSubscribeFailed(t *testing.T) {
	data := make(chan *metric.Metric, 3)
	data <- foo
	data <- bar
	data <- baz

	coll := make([]*metric.Metric, 0, 3)
	e := subscribe(data, func(s *metric.Metric) error {
		coll = append(coll, s)
		return errors.New("fake error")
	})

	var wg sync.WaitGroup
	poller := func() {
		e.Start()
		for i := 0; i < 3; i++ {
			e.Exec(ctx, voidFunc)
		}
		e.End()
		wg.Done()
	}

	wg.Add(1)
	go poller()
	wg.Wait()
	close(data)

	assert.ElementsMatch(t, coll, []*metric.Metric{foo, bar, baz})
}
