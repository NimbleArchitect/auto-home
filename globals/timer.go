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

func (t *timer) Cancel() {
	t.reset <- false
}

// startAndWait starts the timer and waits for the timers to expire
//
// returns: true if we are ending due to the timer and false if we was cancelled
func (t *timer) startAndWait() bool {

	tick := time.NewTimer(t.duration)
	for {
		select {
		case shouldReset := <-t.reset:
			if shouldReset {
				if !tick.Stop() {
					<-tick.C
				}
				tick.Reset(t.duration)
			} else {
				return false
			}
		case <-tick.C:
			return true

		}
	}
}
