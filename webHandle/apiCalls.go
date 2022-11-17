package webHandle

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	event "server/eventManager"
	"time"

	"github.com/google/uuid"
)

type sessionState struct {
	clientid  string
	actionId  string
	timestamp time.Time
	InUse     bool
	Done      chan bool
}

func (h *Handler) register(req requestInfoBlock) {
	var tmp Generic
	var hub jsonHub
	var device jsonDevice

	err := json.Unmarshal(req.Body, &tmp)
	if err != nil {
		log.Printf("error reading body for %s: %s\n", req.Path, err.Error())
	}

	// TODO: lookup client id from session id
	state, ok := h.session[req.Session]
	if !ok {
		log.Println("error invalid session id")
		return
	}

	if *tmp.Method == "hub" {
		raw, _ := tmp.Data.MarshalJSON()
		json.Unmarshal(raw, &hub)
		log.Println("registration for hub", hub.Name)
		// build device list
		err = h.regHubList(hub, state.clientid)
		if err != nil {
			log.Println("Error:", err)
			req.Response.WriteHeader(http.StatusInternalServerError)
			writeFlush(req.Response, "Error: unable to add hub")
			return
		}
	}

	if *tmp.Method == "device" {
		raw, _ := tmp.Data.MarshalJSON()
		json.Unmarshal(raw, &device)
		log.Printf("registration for device \"%s\" (id: %s)", device.Name, device.Id)

		err := h.regDeviceList(device, state.clientid)
		if err != nil {
			log.Println("Error:", err)
			req.Response.WriteHeader(http.StatusInternalServerError)
			writeFlush(req.Response, "Error: unable to add device")
			return
		}
	}

	// h.addDeviceActionList(newUuid)
	log.Println("registration successful")
	writeFlush(req.Response, `{"result": {"status":"ok","msg":""}}\n`)

}

// processEvent reads the incoming event as json, checks it valid and passes it off to the event manager
func (h *Handler) processEvent(req requestInfoBlock, clientId string) {
	var jEvent JsonEvent

	// scanner := bufio.NewScanner(r.Body)
	// for scanner.Scan() {
	// 	ln := scanner.Text()
	// 	log.Println("scanner recieved", ln)
	// }

	fmt.Println(">> jsonMessage:", string(req.Body))
	// err := json.NewDecoder(req.Request.Body).Decode(&jEvent)
	err := json.Unmarshal(req.Body, &jEvent)
	if err != nil {
		log.Println(err)
	}

	switch jEvent.Method {
	case "event":
		// make sure the client id has ownership of the device id
		// and verify device exists
		if !h.HomeManager.DeviceExistsWithClientId(jEvent.Data.Id, clientId) {
			log.Println("Error: invalid device id")
			req.Response.WriteHeader(http.StatusBadRequest)
			writeFlush(req.Response, "Error: incorrect device id")
			return
		}
	}

	// TODO: validate the timestamp
	// jEvent.Data.Timestamp

	msg := event.EventMsg{
		Id: jEvent.Data.Id,
		// Timestamp:  time.Now(),
		Properties: jEvent.Data.Properties,
	}

	// add event to queue
	h.EventManager.AddEvent(msg)

}

func (h *Handler) callV1api(req requestInfoBlock) {
	// is the user logged in
	if !h.isConnected(req) {
		// not logged in
		switch req.Components[1] {
		case "connect":
			fmt.Println("/connect") // api login
			if req.Request.Method == "POST" {
				h.doLogin(req)
			}
		default:
			log.Println("not logged in")
			writeFlush(req.Response, "not logged in")
			return
		}

	} else {
		// all ok
		switch req.Components[1] {
		case "register":
			fmt.Println("/register")
			h.register(req)

		case "actions":
			fmt.Println("/actions")
			if client, ok := h.doActions(req); ok {
				h.HomeManager.SetClient(client.clientid, req.Response, req.Request)
				ctx := req.Request.Context()
				select {
				case <-ctx.Done():
				case <-client.Done:
				}
				log.Println("finished /actions")
			}

		case "event":
			fmt.Println("/event")
			val, ok := h.session[req.Session]
			if !ok {
				log.Println("invalid session")
				return
			}
			h.processEvent(req, val.clientid)

		default:
			log.Println("unknown url:", req.Path)
		}

	}

	// time.Sleep(4 * time.Second)
}

func (h *Handler) isConnected(req requestInfoBlock) bool {
	log.Println("header:", req.Request.Header)

	if len(req.Session) <= 0 {
		return false
	}

	state, ok := h.session[req.Session]
	if ok {
		if state.timestamp.After(time.Now()) {
			return true
		} else {
			return false
		}
	}

	return false
}

func (h *Handler) doLogin(req requestInfoBlock) bool {
	type userLogin struct {
		User string
		Pass string
	}
	var login userLogin

	now := time.Now()

	rawMsg, err := req.JsonMessage.Data.MarshalJSON()
	if err != nil {
		log.Println("unable to retrieve bytes from generic:", err)
		return false
	}

	err = json.Unmarshal(rawMsg, &login)
	if err != nil {
		log.Println("unable to convert json string", err)
		return false
	}

	val, ok := h.userInfo[login.User]
	if !ok {
		log.Println("username not found")
		return false
	}

	if login.Pass != val.AuthKey {
		log.Println("invalid authKey")
		return false
	}

	// if len(val.RefreshToken) == 0 {
	// 	val.RefreshToken = uuid.New().String()
	// 	val.RefreshTime = now.Add(24 * time.Hour)
	// } else if val.RefreshTime.After(now) {
	// 	// long lived token has expired
	// 	// do I need to do anything
	// }

	// user login is ok so generate session id and and action id then attach clientid
	newSessionID := uuid.New().String()
	newActionID := uuid.New().String()
	h.session[newSessionID] = sessionState{
		clientid:  login.User,
		actionId:  newActionID,
		timestamp: now.Add(24 * 60 * time.Minute),
		InUse:     true,
		Done:      make(chan bool),
	}

	req.Response.Header().Set("session", newSessionID)
	req.Response.WriteHeader(200)
	writeFlush(req.Response,
		fmt.Sprintf(
			`{"result":{"status":"ok","msg":""},"data":{"session":"%s","actionid":"%s"}}`,
			newSessionID, newActionID,
		),
	)

	return true
}

func (h *Handler) doActions(req requestInfoBlock) (sessionState, bool) {
	val, ok := h.session[req.Session]
	if !ok {
		log.Println("invalid session")
		return sessionState{}, false
	}

	actionid := req.Components[2]
	if len(actionid) == 0 {
		log.Println("empty action id")
		return sessionState{}, false
	}

	if val.actionId != actionid {
		log.Println("invalid action id")
		return sessionState{}, false
	}

	// fmt.Println("go go go:", val.clientid)
	// TODO: need to convert back to http3 and check that the waits work correctly and messages are
	//  passed back and forth as needed
	req.Response.WriteHeader(http.StatusAccepted)
	writeFlush(req.Response, "")

	return val, true
}
