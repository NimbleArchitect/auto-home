package globals

import (
	"sync"
	"time"
)

type timer struct {
	lock     sync.Mutex
	reset    chan bool
	duration time.Duration
}

func (t *timer) Reset() {
	t.reset <- true
}

func (t *timer) startAndWait() {

	tick := time.NewTimer(t.duration)
	for {
		select {
		case <-t.reset:
			if !tick.Stop() {
				<-tick.C
			}
			tick.Reset(t.duration)
		case <-tick.C:
			return
		}
	}
}
