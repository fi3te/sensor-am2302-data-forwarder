package control

import "time"

type StatefulTicker struct {
	Ticker          *time.Ticker
	InitialInterval time.Duration
	CurrentInterval time.Duration
}

func (st *StatefulTicker) Reset(interval time.Duration) {
	st.CurrentInterval = interval
	st.Ticker.Reset(interval)
}

func (st *StatefulTicker) ResetToInitialInterval() {
	st.Reset(st.InitialInterval)
}
