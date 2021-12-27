package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type State int

var _ fmt.Stringer = (*State)(nil)

func (s State) String() string {
	switch s {
	case Idle:
		return "idle"
	case Running:
		return "running"
	case Stopped:
		return "stopped"
	}
	return "undefined"
}

const (
	Idle State = iota
	Running
	Stopped
)

type Task struct {
	state atomic.Value

	wg     *sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc

	job      func(*Task)
	tearDown func()
}

func (t *Task) NewTask(job func(*Task)) *Task {
	n := Task{
		wg:  t.wg,
		job: job,
	}
	n.state.Store(Idle)
	n.ctx, n.cancel = context.WithCancel(t.ctx)
	return &n
}

func NewApp() *Task {
	var wg sync.WaitGroup
	task := Task{
		wg:       &wg,
		tearDown: wg.Wait,
	}
	task.state.Store(Idle)

	task.ctx, task.cancel = context.WithCancel(context.Background())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	task.job = func(t *Task) {
		defer t.Stop()
		log.Printf("%v signal received", <-signals)
	}
	return &task
}

func (t *Task) State() State {
	return t.state.Load().(State)
}

func (t *Task) Run(onStop func()) {
	t.setRunning()

	t.tearDown = func() {
		defer t.wg.Done()
		defer t.cancel()
		if onStop != nil {
			onStop()
		}
	}

	t.wg.Add(1)
	go func() {
		defer t.Stop()
		t.job(t)
	}()
}

func (t *Task) Immediate() {
	t.setRunning()

	defer t.Stop()
	t.tearDown = t.cancel
	t.job(t)
}

func (t *Task) Periodic(period time.Duration, onStop func()) {
	t.setRunning()

	ticker := time.NewTicker(period)

	t.tearDown = func() {
		defer t.wg.Done()
		defer t.cancel()
		ticker.Stop()
		if onStop != nil {
			onStop()
		}
	}

	done := t.ctx.Done()

	t.wg.Add(1)
	go func() {
		for t.State() == Running {
			select {
			case <-ticker.C:
				t.job(t)
			case <-done:
				t.Stop()
			}
		}
	}()
}

func (t *Task) Stop() {
	if t.state.CompareAndSwap(Running, Stopped) && t.tearDown != nil {
		t.tearDown()
	}
}

func (t *Task) setRunning() {
	if !t.state.CompareAndSwap(Idle, Running) {
		panic("cant run non idle task")
	}
}
