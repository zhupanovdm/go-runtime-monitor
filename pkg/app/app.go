// Package app contains commonly used application helper functions and types.
package app

import (
	"os"
	"os/signal"
	"syscall"
)

type BuildInfoString string

// TerminationSignal returns initialised channel that will be notified with OS signal when application receives termination sig.
func TerminationSignal() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	return ch
}

func (s BuildInfoString) String() string {
	if s == "" {
		return "N/A"
	}
	return string(s)
}
