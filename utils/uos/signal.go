package uos

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitSignal() {
	var event = make(chan os.Signal)
	signal.Notify(
		event,
		syscall.SIGTERM,
		syscall.SIGABRT,
		syscall.SIGQUIT,
	)

	select {
	case <-event:
		break
	}
}
