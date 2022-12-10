package control

import (
	"os"
	"os/signal"
)

func WaitForInterrupt() {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)
	<-interruptChan
}
