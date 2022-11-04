package globals

import (
	"sync"
	"time"
)

type Global struct {
	lock   sync.Mutex
	vars   map[string]interface{} // varibles that can be referenced across VMs
	timers map[string]*timer      //set reset timers that call a func when they expire
}

func New() *Global {
	return &Global{
		lock:   sync.Mutex{},
		vars:   make(map[string]interface{}),
		timers: make(map[string]*timer),
	}
}

// SetTimer creates a timer that calls the specified function when the timer expires,
//
//	calling SetTime multipkle times will reset the countdown, once set the duration cant be changed
//	if mSec is set to zero the timer is stopped without firing and deleted
func (g *Global) SetTimer(name string, mSec int, call func()) {
	g.lock.Lock()
	val, ok := g.timers[name]
	g.lock.Unlock()

	if ok {
		if mSec == 0 {
			val.Cancel()

			g.lock.Lock()
			delete(g.timers, name)
			g.lock.Unlock()
		} else {
			// exists - we have a timer setup and running, so call the reset
			val.Reset()
		}
	} else {
		// not found - create a new one
		newTimer := timer{
			lock:     sync.Mutex{},
			reset:    make(chan bool, 1),
			duration: time.Duration(mSec) * time.Millisecond,
		}

		g.lock.Lock()
		g.timers[name] = &newTimer
		g.lock.Unlock()
		go func() {
			// then start it
			doCall := newTimer.startAndWait()
			if doCall {
				call()
			}
			g.lock.Lock()
			delete(g.timers, name)
			g.lock.Unlock()
		}()
	}
}
