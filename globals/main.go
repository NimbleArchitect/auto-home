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
// calling SetTime multiple times will reset the countdown, once set the duration cant be changed
// if sec is set to zero the timer is stopped without firing and deleted
// call accepts true if it should run as successfull and false if it was cancelled
func (g *Global) SetTimer(name string, sec float64, call func(bool)) {
	g.lock.Lock()
	val, ok := g.timers[name]
	g.lock.Unlock()

	if ok {
		if sec == 0 {
			val.Cancel()

			g.lock.Lock()
			delete(g.timers, name)
			g.lock.Unlock()
			// we must run call so it kicks back and unlocks the vm
			call(false)
		} else {
			// exists - we have a timer setup and running, so call the reset
			val.Reset()
		}
	} else {
		// not found - create a new one
		newTimer := timer{
			lock:     sync.Mutex{},
			reset:    make(chan bool, 1),
			duration: time.Duration(sec*1000) * time.Millisecond,
		}

		g.lock.Lock()
		g.timers[name] = &newTimer
		g.lock.Unlock()
		go func() {
			// then start it
			doCall := newTimer.startAndWait()
			if doCall {
				call(true)
			}
			g.lock.Lock()
			delete(g.timers, name)
			g.lock.Unlock()
		}()
	}
}

func (g *Global) SetVariable(name string, value interface{}) {
	g.lock.Lock()
	_, ok := g.vars[name]
	if ok {
		g.vars[name] = value
	}
	g.lock.Unlock()

}

func (g *Global) GetVariable(name string) interface{} {
	g.lock.Lock()
	val, ok := g.vars[name]
	g.lock.Unlock()

	if ok {
		return val
	}

	return nil
}
