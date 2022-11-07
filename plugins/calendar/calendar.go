package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type calendarEvent struct {
	nextEvent *calendarEvent
	date      time.Time
	data      interface{}
}

type calendar struct {
	lock       sync.Mutex
	eventStart *calendarEvent
	eventEnd   *calendarEvent
	nextEvent  *calendarEvent
	isRunning  bool
	plugin     *plugin
}

func (c *calendar) start() {
	var now time.Time

	fmt.Println("Start")
	c.isRunning = true

	go func() {
		defer func() { c.isRunning = false }()
		lastKnownTime := time.Now()
		c.nextEvent = c.eventStart.nextEvent

		for {
			next := c.latest()
			if next == nil {
				break
			}

			now = time.Now()
			if lastKnownTime.Sub(now).Seconds() > 5 {
				c.eventReCheck()
			}
			lastKnownTime = now
			if next.date.Sub(now).Seconds() <= 0 {
				go c.fireEvent(next.data)
				c.loadNextEvent()
			}
			time.Sleep(1 * time.Second)
		}

		c.eventStart = nil
		c.nextEvent = nil
		c.eventEnd = nil
		c.isRunning = false
	}()

}

func (c *calendar) loadNextEvent() {
	c.lock.Lock()
	c.nextEvent = c.nextEvent.nextEvent
	c.lock.Unlock()
}

func (c *calendar) eventReCheck() {
	fmt.Print("FIXME: recheck next timer")
	// TODO: we might have missed events as the system time was changed
	// send a trigger back to the server so the users can attach an onchange event
	// also fast-forward over any missed events so we dont sent out loads of alerts

}

func (c *calendar) add(t time.Time, event interface{}) error {
	if t.Before(time.Now()) {
		return errors.New("unable to add dates in the past")
	}

	evt := calendarEvent{
		date: t,
		data: event,
	}

	c.lock.Lock()
	if c.eventStart == nil {
		c.eventEnd = &calendarEvent{}
		evt.nextEvent = c.eventEnd
		c.eventStart = &calendarEvent{nextEvent: &evt}
		c.nextEvent = c.eventStart
	} else {
		this := c.eventStart
		for {
			if this.nextEvent.date.After(t) || this.nextEvent.nextEvent == nil {
				evt.nextEvent = this.nextEvent
				this.nextEvent = &evt
				break
			}
			this = this.nextEvent
		}
		c.nextEvent = c.eventStart
	}

	c.lock.Unlock()
	if !c.isRunning {
		c.start()
	}

	return nil
}

func (c *calendar) latest() *calendarEvent {
	c.lock.Lock()
	out := c.nextEvent
	if out.nextEvent == nil {
		out = nil
	}
	c.lock.Unlock()
	return out
}
