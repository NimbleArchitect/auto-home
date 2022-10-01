package event

import (
	"log"
)

// event queue reads the event coming in on the channel
// processes the event into its area,
// lines out the message on the correct outbound channel

type Manager struct {
	events []EventMsg
	//list of historic events, number of event held shoud be based on length of time
	eventHistory      []EventMsg
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
	Timestamp  string
}

// func loopEvents() {
// 	events := make([]string, 200)
// 	addEvent := make(chan string)
// 	eventId := make(chan int, 50)
// 	removeEvent := make(chan int, 50)
// 	closeEventManager := make(chan bool)

// 	go eventLoop(events, eventId, removeEvent)
// 	go eventManager(events, addEvent, eventId, removeEvent, closeEventManager)

// 	go func() {
// 		// time.Sleep(time.Duration(rand.Intn(4000)) * time.Millisecond)
// 		for i := 0; i < 30; i++ {
// 			// create 20 events
// 			fmt.Println("new event a" + strconv.Itoa(i))
// 			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
// 			addEvent <- "event a" + strconv.Itoa(i)
// 		}
// 	}()

// 	<-closeEventManager
// }

func NewManager(eventQueueLen int, bufferLen int) *Manager {
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

func (e *Manager) EventManager() {
	var eventCount, headPos int

	log.Println("starting EventManager")

	for {
		select {
		case msg := <-e.chAdd:
			log.Printf("add event (%d/%d/%d): %+v\n", headPos, eventCount, len(e.events), msg)

			if eventCount < len(e.events) {
				e.events[headPos] = msg
				eventCount += 1
				e.chCurrentEvent <- headPos
				headPos += 1
				if headPos >= len(e.events) {
					headPos = 0
				}
			} else {
				log.Println("too many events to process")
			}

		case <-e.chRemove:
			log.Printf("remove event (%d/%d/%d)\n", headPos, eventCount, len(e.events))
			eventCount -= 1

		case <-e.closeEventManager:
			return
		}
	}
}

func (e *Manager) AddEvent(event EventMsg) {
	e.chAdd <- event
}

type EventLoop interface {
	Trigger(string, string, []map[string]interface{}) error
	// SaveState() (interface{}, error)
}

func (e *Manager) EventLoop(looper EventLoop) {
	log.Println("starting EventLoop")

	for {
		select {
		case evtid := <-e.chCurrentEvent:
			msg := e.events[evtid]

			log.Println("processing event", msg.Id)

			// this is where we actually do something with the event
			//first we save the system state
			// state, err := looper.SaveState()
			// if err != nil {
			// 	log.Println("unable to save state:", err)
			// } else {
			// then if everything went ok we call the event trigger function and
			//  pass in message properties and the saved state
			go func() { //go call so we can run the trigger concurrently
				err := looper.Trigger(msg.Id, msg.Timestamp, msg.Properties)
				if err != nil {
					log.Println("event error", err)
				}
				// }

				// signal to the remove channel that we have finished processing the event
				e.chRemove <- evtid
			}()
		case <-e.closeEventLoop:
			return
		}
	}
}
