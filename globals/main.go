package globals

import (
	"fmt"
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

type timer struct {
	lock     sync.Mutex
	reset    chan bool
	duration time.Duration
}

func (t *timer) Reset() {
	t.reset <- true
}

func (t *timer) start() {

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

func (g *Global) SetTimer(name string, mSec int, call func()) {
	g.lock.Lock()
	val, ok := g.timers[name]
	g.lock.Unlock()

	if ok {
		// exists - we have a timer setup and running, so call the reset
		val.Reset()
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
		// then start it
		go func() {
			newTimer.start()
			fmt.Println("setTime start returned")
			call()
			g.lock.Lock()
			delete(g.timers, name)
			g.lock.Unlock()
		}()
	}
}
