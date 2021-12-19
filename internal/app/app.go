package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var Context context.Context
var Cancel context.CancelFunc
var WG = &sync.WaitGroup{}

var termSignal chan os.Signal

func init() {
	termSignal = make(chan os.Signal, 1)
	signal.Notify(termSignal, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	Context, Cancel = context.WithCancel(context.Background())
}

func Periodic(period time.Duration, e Executor) {
	ticker := time.NewTicker(period)
	ctx, cancel := context.WithCancel(Context)

	e.Start()
	WG.Add(1)
	go func(done <-chan struct{}) {
		defer WG.Done()
		defer cancel()
		defer e.End()

	loop:
		for {
			select {
			case <-ticker.C:
				e.Exec(ctx, cancel)
			case <-done:
				ticker.Stop()
				break loop
			}
		}
	}(ctx.Done())
}

func Serve() {
	s := <-termSignal
	Cancel()
	WG.Wait()
	log.Printf("%v signal received", s)
}

type Executor interface {
	Start()
	Exec(ctx context.Context, cancel context.CancelFunc)
	End()
}

type ExecutorHandler struct {
	OnStart func()
	OnExec  func(ctx context.Context, cancel context.CancelFunc)
	OnEnd   func()
}

var _ Executor = (*ExecutorHandler)(nil)

func (e ExecutorHandler) Start() {
	if e.OnStart != nil {
		e.OnStart()
	}
}

func (e ExecutorHandler) Exec(ctx context.Context, cancel context.CancelFunc) {
	e.OnExec(ctx, cancel)
}

func (e ExecutorHandler) End() {
	if e.OnEnd != nil {
		e.OnEnd()
	}
}
