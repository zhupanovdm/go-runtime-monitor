package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

var ctx = context.TODO()
var voidFunc = func() {}

func TestPublish(t *testing.T) {
	data := make(chan string)
	e := publish(data, func(c chan<- string) {
		c <- "foo"
		c <- "bar"
		c <- "baz"
	})

	go func() {
		e.Start()
		e.Exec(ctx, voidFunc)
		e.End()
	}()

	assert.Equal(t, <-data, "foo")
	assert.Equal(t, <-data, "bar")
	assert.Equal(t, <-data, "baz")

	_, ok := <-data
	assert.False(t, ok)
}

func TestSubscribe(t *testing.T) {
	data := make(chan string, 3)
	data <- "foo"
	data <- "bar"
	data <- "baz"

	coll := make([]string, 0, 3)
	e := subscribe(data, func(s string) error {
		coll = append(coll, s)
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

	assert.ElementsMatch(t, coll, []string{"foo", "bar", "baz"})
}

func TestSubscribePartial(t *testing.T) {
	data := make(chan string, 3)
	data <- "foo"
	data <- "bar"

	coll := make([]string, 0, 3)
	e := subscribe(data, func(s string) error {
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
	assert.ElementsMatch(t, coll, []string{"foo", "bar"})

	data <- "baz"
	wg.Add(1)
	go poller()
	wg.Wait()
	close(data)

	assert.ElementsMatch(t, coll, []string{"foo", "bar", "baz"})
}

func TestSubscribeClosed(t *testing.T) {
	data := make(chan string, 3)
	data <- "foo"
	data <- "bar"

	coll := make([]string, 0, 3)
	e := subscribe(data, func(s string) error {
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
	assert.ElementsMatch(t, coll, []string{"foo", "bar"})

	close(data)
	wg.Add(1)
	go poller()
	wg.Wait()

	assert.ElementsMatch(t, coll, []string{"foo", "bar"})
}

func TestSubscribeFailed(t *testing.T) {
	data := make(chan string, 3)
	data <- "foo"
	data <- "bar"
	data <- "baz"

	coll := make([]string, 0, 3)
	e := subscribe(data, func(s string) error {
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

	assert.ElementsMatch(t, coll, []string{"foo", "bar", "baz"})
}
