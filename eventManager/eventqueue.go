package event

import (
	"server/logger"
	"time"
)

var debugLevel int

// event queue reads the event coming in on the channel
// processes the event into its area,
// lines out the message on the correct outbound channel

type Manager struct {
	events            []EventMsg
	chAdd             chan EventMsg
	chCurrentEvent    chan int
	chRemove          chan int
	closeEventManager chan bool
	closeEventLoop    chan bool
}

// type evt struct {
// 	id   int
// 	read int
// }

type EventMsg struct {
	Id         string
	EventId    string
	Properties []map[string]interface{}
	Timestamp  time.Time
}

// NewManager
//
//	eventQueueLen length of the events array
//	bufferLen number of events that can be spawned for concurrent processing
func NewManager(eventQueueLen int, bufferLen int) *Manager {
	debugLevel = logger.GetDebugLevel()

	m := Manager{
		events:            make([]EventMsg, eventQueueLen),
		chAdd:             make(chan EventMsg),
		chCurrentEvent:    make(chan int, bufferLen),
		chRemove:          make(chan int, bufferLen),
		closeEventManager: make(chan bool),
		closeEventLoop:    make(chan bool),
	}
	return &m
}

func (e *Manager) Shutdown() {
	log := logger.New("eventManager.Shutdown", &debugLevel)
	log.Trace("start")
	e.closeEventManager <- true
	e.closeEventLoop <- true
	log.Trace("stop")
}

func (e *Manager) EventManager() {
	var eventCount, headPos int
	log := logger.New("EventManager", &debugLevel)
	log.Info("starting EventManager")

	for {
		select {
		case msg := <-e.chAdd:
			log.Infof("add event (%d/%d/%d): %+v\n", headPos, eventCount, len(e.events), msg)

			if eventCount < len(e.events) {
				e.events[headPos] = msg
				eventCount += 1
				e.chCurrentEvent <- headPos
				headPos += 1
				if headPos >= len(e.events) {
					headPos = 0
				}
			} else {
				log.Error("too many events to process")
			}

		case <-e.chRemove:
			log.Infof("remove event (%d/%d/%d)\n", headPos, eventCount, len(e.events))
			eventCount -= 1

		case <-e.closeEventManager:
			log.Info("stopping EventManager")
			return
		}
	}
}

func (e *Manager) AddEvent(event EventMsg) {
	// we reset the timestamp to the point it enters the queue
	event.Timestamp = time.Now()
	e.chAdd <- event
}

type EventLoop interface {
	Trigger(int, string, time.Time, []map[string]interface{}) error
	// SaveState() (interface{}, error)
}

func (e *Manager) EventLoop(looper EventLoop) {
	log := logger.New("EventLoop", &debugLevel)

	log.Info("starting EventLoop")
	loopId := 0

	for {
		select {
		case evtid := <-e.chCurrentEvent:
			msg := e.events[evtid]

			log.Info("processing event", msg.Id)

			// this is where we actually do something with the event
			//  we call the event trigger function and
			//  pass in message properties and the saved state
			//go call so we can run the trigger concurrently
			go func() {
				id := loopId

				err := looper.Trigger(id, msg.Id, msg.Timestamp, msg.Properties)
				if err != nil {
					log.Error("event error", err)
				}
			}()

			// signal to the remove channel that we have finished processing the event
			e.chRemove <- evtid

		case <-e.closeEventLoop:
			log.Info("stopping EventLoop")
			return
		}
		loopId++
	}
}
