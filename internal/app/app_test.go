package app

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecutorHandlerEnd(t *testing.T) {
	var executed bool
	tests := []struct {
		name         string
		sample       Executor
		wantExecuted bool
	}{
		{
			name: "Basic test",
			sample: &ExecutorHandler{
				OnEnd: func() {
					executed = true
				},
			},
			wantExecuted: true,
		},
		{
			name:   "Empty handler",
			sample: &ExecutorHandler{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executed = false
			assert.NotPanics(t, func() {
				tt.sample.End()
			})
			if tt.wantExecuted {
				assert.True(t, executed, "handler not executed")
			} else {
				assert.False(t, executed, "handler executed")
			}
		})
	}
}

func TestExecutorHandlerExec(t *testing.T) {
	var executed bool
	tests := []struct {
		name         string
		sample       Executor
		wantExecuted bool
		wantPanic    bool
	}{
		{
			name: "Basic test",
			sample: &ExecutorHandler{
				OnExec: func(context.Context, context.CancelFunc) {
					executed = true
				},
			},
			wantExecuted: true,
		},
		{
			name:      "Empty handler",
			sample:    &ExecutorHandler{},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executed = false
			f := func() {
				tt.sample.Exec(context.TODO(), func() {})
			}
			if tt.wantPanic {
				assert.Panics(t, f)
			} else {
				assert.NotPanics(t, f)
			}
			if tt.wantExecuted {
				assert.True(t, executed, "handler not executed")
			} else {
				assert.False(t, executed, "handler executed")
			}
		})
	}
}

func TestExecutorHandlerStart(t *testing.T) {
	var executed bool
	tests := []struct {
		name         string
		sample       Executor
		wantExecuted bool
	}{
		{
			name: "Basic test",
			sample: &ExecutorHandler{
				OnStart: func() {
					executed = true
				},
			},
			wantExecuted: true,
		},
		{
			name:   "Empty handler",
			sample: &ExecutorHandler{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executed = false
			assert.NotPanics(t, func() {
				tt.sample.Start()
			})
			if tt.wantExecuted {
				assert.True(t, executed, "handler not executed")
			} else {
				assert.False(t, executed, "handler executed")
			}
		})
	}
}
