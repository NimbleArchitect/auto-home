package webHandle

import (
	"encoding/json"
	"fmt"

	"net/http"
	event "server/eventManager"
	log "server/logger"
	"time"

	"github.com/google/uuid"
)

type sessionState struct {
	clientid  string // id of the connecting client
	actionId  string // new generated action id
	timestamp time.Time
	InUse     bool      // is this session state currently in use and registered
	Done      chan bool // chennel recieves true when this session has completed
}

type userLogin struct {
	ReturnId bool   `json:"returnid"` // if true we want the action id and the event id to be returned
	User     string `json:"user"`
	Pass     string `json:"pass"`
}

func (h *Handler) register(req requestInfoBlock) {
	var tmp Generic
	var hub jsonHub
	var device jsonDevice

	err := json.Unmarshal(req.Body, &tmp)
	if err != nil {
		log.Errorf("error reading body for %s: %s\n", req.Path, err.Error())
	}

	// TODO: lookup client id from session id
	state, ok := h.session[req.Session]
	if !ok {
		log.Error("error invalid session id")
		return
	}

	if *tmp.Method == "hub" {
		raw, _ := tmp.Data.MarshalJSON()
		json.Unmarshal(raw, &hub)
		log.Info("registration for hub", hub.Name)
		// build device list
		err = h.regHubList(hub, state.clientid)
		if err != nil {
			log.Error(err)
			req.Response.WriteHeader(http.StatusInternalServerError)
			writeFlush(req.Response, "Error: unable to add hub")
			return
		}
	}

	if *tmp.Method == "device" {
		raw, _ := tmp.Data.MarshalJSON()
		json.Unmarshal(raw, &device)
		log.Infof("registration for device \"%s\" (id: %s)", device.Name, device.Id)

		err := h.regDeviceList(device, state.clientid)
		if err != nil {
			log.Error(err)
			req.Response.WriteHeader(http.StatusInternalServerError)
			writeFlush(req.Response, "Error: unable to add device")
			return
		}
	}

	// h.addDeviceActionList(newUuid)
	log.Info("registration successful")
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

	log.Info(">> jsonMessage:", string(req.Body))
	// err := json.NewDecoder(req.Request.Body).Decode(&jEvent)
	err := json.Unmarshal(req.Body, &jEvent)
	if err != nil {
		log.Error(err)
	}

	switch jEvent.Method {
	case "event":
		// make sure the client id has ownership of the device id
		// and verify device exists
		if !h.HomeManager.DeviceExistsWithClientId(jEvent.Data.Id, clientId) {
			log.Error("invalid device id")
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

	log.Debug("req.Components[1] =", req.Components[1])

	// is the user logged in
	if !h.isConnected(req) {
		// not logged in
		switch req.Components[1] {
		case "connect":
			log.Info("/connect") // api login
			if req.Request.Method == "POST" {
				h.doLogin(req)
			}
		default:
			log.Info("not logged in")
			writeFlush(req.Response, "not logged in")
			return
		}

	} else {
		// all ok
		switch req.Components[1] {
		case "connect":
			log.Info("/connect") // api login
			if req.Request.Method == "POST" {
				h.doLogin(req)
			}

		case "register":
			log.Info("/register")
			h.register(req)

		case "actions":
			log.Info("/actions")
			if client, ok := h.doActions(req); ok {
				h.HomeManager.SetClient(client.clientid, req.Response, req.Request)
				ctx := req.Request.Context()
				select {
				case <-ctx.Done():
				case <-client.Done:
				}
				log.Info("finished /actions")
			}

		case "event":
			log.Info("/event")
			val, ok := h.session[req.Session]
			if !ok {
				log.Error("invalid session")
				return
			}
			h.processEvent(req, val.clientid)

		case "device":
			log.Info("/device")
			if req.Request.Method == "GET" {
				switch len(req.Components) {
				case 2: // v1/device
					bytesout := h.HomeManager.WebAllDeviceAsJson()
					req.Response.WriteHeader(http.StatusOK)
					writeFlush(req.Response, string(bytesout))
				case 3: // v1/device/deviceid
					bytesout := h.HomeManager.WebDeviceAsJson(req.Components[2])
					req.Response.WriteHeader(http.StatusOK)
					writeFlush(req.Response, string(bytesout))
				case 4: // v1/device/deviceid/propertyname
					if len(req.Query) == 0 { // no options provided so we return the property state
						bytesout := h.HomeManager.WebDevicePropertyAsJson(req.Components[2], req.Components[3])
						req.Response.WriteHeader(http.StatusOK)
						writeFlush(req.Response, string(bytesout))
					} else { // if we have options then this coud be a set request
						if value := req.Query.Get("setstate"); len(value) > 0 {
							ok := h.HomeManager.WebSetDeviceProperty(req.Components[2], req.Components[3], value)
							if ok {
								req.Response.WriteHeader(http.StatusOK)
								writeFlush(req.Response, "")
							}
						} else if percent := req.Query.Get("setpercent"); len(percent) > 0 {
							ok := h.HomeManager.WebSetDevicePropertyPercent(req.Components[2], req.Components[3], percent)
							if ok {
								req.Response.WriteHeader(http.StatusOK)
								writeFlush(req.Response, "")
							}
						} else { // return error
							req.Response.WriteHeader(http.StatusBadRequest)
							writeFlush(req.Response, "")
						}

					}
				}

				// } else if req.Request.Method == "PATCH" {

			}

		default:
			log.Error("unknown url:", req.Path)
		}

	}

	// time.Sleep(4 * time.Second)
}

func (h *Handler) isConnected(req requestInfoBlock) bool {

	log.Debug("header:", req.Request.Header)

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
	var login userLogin
	var newActionID string
	var jsonResponse string

	now := time.Now()

	rawMsg, err := req.JsonMessage.Data.MarshalJSON()
	if err != nil {
		log.Error("unable to retrieve bytes from generic:", err)
		return false
	}

	err = json.Unmarshal(rawMsg, &login)
	if err != nil {
		log.Error("unable to convert json string", err)
		return false
	}

	val, ok := h.userInfo[login.User]
	if !ok {
		log.Error("username not found")
		return false
	}

	if login.Pass != val.AuthKey {
		log.Error("invalid authKey")
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

	// are we expected to return an action id along with the session id
	if login.ReturnId {
		newActionID = uuid.New().String()
		jsonResponse = fmt.Sprintf(
			`{"result":{"status":"ok","msg":""},"data":{"session":"%s","actionid":"%s"}}`,
			newSessionID, newActionID,
		)
	} else {
		jsonResponse = fmt.Sprintf(
			`{"result":{"status":"ok","msg":""},"data":{"session":"%s"}}`,
			newSessionID,
		)
	}

	h.session[newSessionID] = sessionState{
		clientid:  login.User,
		actionId:  newActionID,
		timestamp: now.Add(24 * 60 * time.Minute),
		InUse:     true,
		Done:      make(chan bool, 2), // TODO: setting this to 1 causes a random lock up during shutdown need to work out why
	}

	req.Response.Header().Set("session", newSessionID)
	req.Response.WriteHeader(200)
	writeFlush(req.Response, jsonResponse)

	return true
}

// doAction checks the sessionid and actionid in the requestInfoBlock and returns a matching sessionState and a success bool
func (h *Handler) doActions(req requestInfoBlock) (sessionState, bool) {

	val, ok := h.session[req.Session]
	if !ok {
		log.Error("invalid session")
		return sessionState{}, false
	}

	actionid := req.Components[2]
	if len(actionid) == 0 {
		log.Error("empty action id")
		return sessionState{}, false
	}

	if val.actionId != actionid {
		log.Error("invalid action id")
		return sessionState{}, false
	}

	// fmt.Println("go go go:", val.clientid)
	// TODO: need to convert back to http3 and check that the waits work correctly and messages are
	//  passed back and forth as needed
	req.Response.WriteHeader(http.StatusAccepted)
	writeFlush(req.Response, "")

	return val, true
}
