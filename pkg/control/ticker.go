package control

import (
	"time"

	"github.com/fi3te/sensor-am2302-data-forwarder/pkg/logging"
)

func RunTicker(interval time.Duration, stopTickerChan chan bool, callback func(*StatefulTicker, time.Time, *logging.LogSetup), ls *logging.LogSetup) {
	ticker := time.NewTicker(interval)
	statefulTicker := StatefulTicker{ticker, interval, interval}

	callback(&statefulTicker, time.Now(), ls)
	for {
		select {
		case t := <-ticker.C:
			callback(&statefulTicker, t, ls)
		case <-stopTickerChan:
			return
		}
	}
}
