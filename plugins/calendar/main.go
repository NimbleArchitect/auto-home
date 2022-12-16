package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

const (
	FLAG_NOTSET = iota
	FLAG_MINUTE
	FLAG_HOUR
	FLAG_DAY
	FLAG_WEEK
	FLAG_MONTH
	FLAG_YEAR
)

type Event struct {
	Id           string
	Created      time.Time
	NextTrigger  time.Time
	firstTrigger time.Time
	CreatedBy    string
	Notify       []string // users/groups to notify when this event fires
	Msg          string
	Location     string
	RepeatCount  int // how many times to repeat
	RepeatEvery  int // repeat evey "RepeatCount" minute/hour/day/week/month/year
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
	// theTime := time.Date(2022, 8, 15, 14, 30, 45, 100, time.Local)
	// utcTime := theTime.UTC()

	evt := Event{
		Id:          "random",
		Created:     time.Now(),
		NextTrigger: time.Date(2022, 8, 15, 14, 30, 45, 100, time.Local),
		CreatedBy:   "me",
		Notify:      []string{"me"},
		Msg:         "go to work",
		Location:    "office",
		RepeatCount: 4,
		RepeatEvery: FLAG_MINUTE,
	}
	c.addEvent(evt)

	evt = Event{
		Id:          "newrand",
		Created:     time.Now(),
		NextTrigger: time.Date(2022, 11, 28, 11, 20, 30, 0, time.Local),
		CreatedBy:   "me",
		Notify:      []string{"me"},
		Msg:         "go home",
		Location:    "home",
		RepeatCount: 1,
		RepeatEvery: FLAG_DAY,
	}
	c.addEvent(evt)

}

func (c *calendar) AddEvent(data []byte) {
	var event Event

	// TODO: need a better way to convert to the Event type auto-magically
	err := json.Unmarshal(data, &event)
	if err != nil {
		return
	}

	c.addEvent(event)
}

func (c *calendar) GetEvent(data []byte) *Event {
	var eventid struct {
		Id string
	}

	err := json.Unmarshal(data, &eventid)
	if err != nil {
		return nil
	}

	return c.getEvent(eventid.Id)
}

func (c *calendar) GetAllEvents(data []byte) []Event {

	return c.getAllEvents()
}

func (c *calendar) GetEventByDate(data []byte) []Event {
	var arguments struct {
		Start time.Time
		End   time.Time
	}

	err := json.Unmarshal(data, &arguments)
	if err != nil {
		fmt.Println("json decode error:", err)
		return nil
	}

	fmt.Println("start:", arguments.Start)
	fmt.Println("end:", arguments.End)
	return c.getEventByDate(arguments.Start, arguments.End)
}

func (c *calendar) GetEventToday(data []byte) []Event {
	n := time.Now()

	start := time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, n.Location())
	end := time.Date(n.Year(), n.Month(), n.Day(), 23, 59, 59, 999999999, n.Location())
	fmt.Println("start:", start)
	fmt.Println("end:", end)
	return c.getEventByDate(start, end)
}

func (c *calendar) addEvent(event Event) {

	event.firstTrigger = event.NextTrigger

	event.updateNextTrigger()

	d := event.NextTrigger

	fmt.Println("adding:", event.NextTrigger)

	if err := c.add(d, event); err != nil {
		fmt.Println("error adding event:", err)
	}
}

func (c *calendar) getEvent(id string) *Event {
	var event Event

	current := c.eventStart
	for {
		current = current.nextEvent
		if current == c.eventEnd {
			break
		}
		event = current.data.(Event)
		if id == event.Id {
			return &event
		}
	}
	return nil
}

func (c *calendar) getEventByDate(start time.Time, end time.Time) []Event {
	var validEvents []Event

	events := c.getAllEvents()

	for _, event := range events {
		if event.RepeatCount > FLAG_NOTSET {
			// calculate if events fall within range
			nextDate := event.firstTrigger
			counter := 0
			for {
				counter++
				if counter > 1000 {
					// TODO: needs a proper error merror message
					fmt.Println("too many events, reduce your time frame")
					break
				}

				newDate := event.calculateNextTrigger(nextDate, start)
				if newDate.After(end) {
					break
				}
				event.NextTrigger = newDate
				validEvents = append(validEvents, event)
				start = newDate
			}
		} else {
			if event.NextTrigger.After(start) && event.NextTrigger.Before(end) {
				validEvents = append(validEvents, event)
			}
		}
	}

	return validEvents
}

func (c *calendar) getAllEvents() []Event {
	var event []Event

	current := c.eventStart
	for {
		current = current.nextEvent
		if current == c.eventEnd {
			break
		}
		event = append(event, current.data.(Event))
	}

	return event
}

func (c *calendar) fireEvent(event interface{}) {
	evt := event.(Event)
	fmt.Println("fire event:", evt.NextTrigger, evt.Msg)

	c.plugin.Call("onevent", evt)

	evt.updateNextTrigger()
	c.add(evt.NextTrigger, evt)
}

func (e *Event) updateNextTrigger() {
	now := time.Now()
	if e.NextTrigger.Before(now) {
		if e.RepeatCount > FLAG_NOTSET {
			nextDate := e.calculateNextTrigger(e.NextTrigger, now)
			e.NextTrigger = nextDate
		}
	}
}

func (e *Event) calculateNextTrigger(nextDate time.Time, now time.Time) time.Time {
	for {
		if nextDate.After(now) {
			break
		}

		switch e.RepeatEvery {
		case FLAG_MINUTE:
			nextDate = nextDate.Add(time.Duration(e.RepeatCount) * time.Minute)
		case FLAG_HOUR:
			nextDate = nextDate.Add(time.Duration(e.RepeatCount) * time.Hour)
		case FLAG_DAY:
			nextDate = nextDate.AddDate(0, 0, e.RepeatCount)
		case FLAG_MONTH:
			nextDate = nextDate.AddDate(0, e.RepeatCount, 0)
		case FLAG_YEAR:
			nextDate = nextDate.AddDate(e.RepeatCount, 0, 0)
		}
	}
	return nextDate
}

// func (e *Event) calculateNextTrigger() {
// 	now := time.Now()
// 	if e.NextTrigger.Before(now) {
// 		if e.RepeatCount > FLAG_NOTSET {
// 			nextDate := e.NextTrigger
// 			for {
// 				if nextDate.After(now) {
// 					break
// 				}

// 				switch e.RepeatEvery {
// 				case FLAG_MINUTE:
// 					nextDate = nextDate.Add(time.Duration(e.RepeatCount) * time.Minute)
// 				case FLAG_HOUR:
// 					nextDate = nextDate.Add(time.Duration(e.RepeatCount) * time.Hour)
// 				case FLAG_DAY:
// 					nextDate = nextDate.AddDate(0, 0, e.RepeatCount)
// 				case FLAG_MONTH:
// 					nextDate = nextDate.AddDate(0, e.RepeatCount, 0)
// 				case FLAG_YEAR:
// 					nextDate = nextDate.AddDate(e.RepeatCount, 0, 0)
// 				}
// 			}
// 			e.NextTrigger = nextDate
// 		}
// 	}
// }
