package control

import (
	"time"
)

func RunTicker(interval time.Duration, stopTickerChan chan bool, callback func(*StatefulTicker, time.Time)) {
	ticker := time.NewTicker(interval)
	statefulTicker := StatefulTicker{ticker, interval, interval}

	callback(&statefulTicker, time.Now())
	for {
		select {
		case t := <-ticker.C:
			callback(&statefulTicker, t)
		case <-stopTickerChan:
			return
		}
	}
}
