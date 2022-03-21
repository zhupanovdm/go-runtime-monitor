// Package app contains commonly used application helper functions and types.
package app

import (
	"os"
	"os/signal"
	"syscall"
)

// TerminationSignal returns initialised channel that will be notified with OS signal when application receives termination sig.
func TerminationSignal() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	return ch
}
