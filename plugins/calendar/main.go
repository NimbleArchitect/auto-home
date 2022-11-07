package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	FLAG_NOTSET = iota
	FLAG_HOUR
	FLAG_DAY
	FLAG_WEEK
	FLAG_MONTH
	FLAG_YEAR
)

type Event struct {
	Id          string
	Created     time.Time
	NextTrigger time.Time
	CreatedBy   string
	Notify      []string // users/groups to notify when this event fires
	Msg         string
	Location    string
	RepeatCount int // how many times to repeat
	RepeatEvery int // repeat evey "RepeatCount" minute/hour/day/week/month/year
}

func main() {
	p := Connect()

	cal := &calendar{
		lock:   sync.Mutex{},
		plugin: p,
	}

	LoadEvents(cal)

	p.Register("calendar", cal)

	err := p.Done()
	if err != nil {
		fmt.Println(err)
	}

}

func LoadEvents(c *calendar) {
	evt := Event{
		Id:          "random",
		Created:     time.Now(),
		NextTrigger: time.Now().Add(time.Second * 5),
		CreatedBy:   "me",
		Notify:      []string{"me"},
		Msg:         "go to work",
		Location:    "office",
		RepeatCount: 1,
		RepeatEvery: FLAG_DAY,
	}

	c.AddEvent(evt)

}

func (c *calendar) AddEvent(event Event) {
	d := event.NextTrigger
	c.Add(d, event)
}

func (c *calendar) fireEvent(event interface{}) {
	evt := event.(Event)
	fmt.Println("fire event:", evt.NextTrigger, evt.Msg)

	c.plugin.Call("onevent", evt)
}
