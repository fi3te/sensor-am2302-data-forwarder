package control

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitForInterrupt() {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)
	<-interruptChan
}
