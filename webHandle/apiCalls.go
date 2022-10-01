package webHandle

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	event "server/eventManager"
)

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var tmp Generic
	var hub jsonHub
	var device jsonDevice
	// var connector jsonConnector
	// var plugin jsonPlugin

	err := json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		log.Printf("error reading body for %s: %s\n", r.URL.Path, err.Error())
	}
	sessionid := r.Header.Get("session")

	// TODO: lookup client id from session id
	clientInfo, ok := h.lookupSessionClient(sessionid)
	if !ok {
		log.Println("error invalid session id")
		return
	}

	if *tmp.Method == "hub" {
		raw, _ := tmp.Data.MarshalJSON()
		json.Unmarshal(raw, &hub)
		log.Println("registration for hub", hub.Name)
		// build device list
		err = h.regHubList(hub, clientInfo.ClientId)
		if err != nil {
			log.Println("Error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error: unable to add hub"))
			return
		}
	}

	if *tmp.Method == "device" {
		raw, _ := tmp.Data.MarshalJSON()
		json.Unmarshal(raw, &device)
		log.Printf("registration for device \"%s\" (id: %s)", device.Name, device.Id)

		err := h.regDeviceList(device, clientInfo.ClientId)
		if err != nil {
			log.Println("Error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error: unable to add device"))
			return
		}
	}

	if *tmp.Method == "listen" {
		// TODO: check session id and generate actionid
		// fmt.Println(">>", r.Header)

		// clientInfo := h.lookupSessionClient(sessionid)
		// raw, _ := tmp.Data.MarshalJSON()
		// json.Unmarshal(raw, &connector)
		h.addDeviceActionList(clientInfo.ClientId, sessionid)
		w.Write([]byte((`{"result": {"status":"ok","msg":""}, "data": {"id": "` + clientInfo.actionId + `"}}\n`)))
	}

	// h.addDeviceActionList(newUuid)
	fmt.Println("success")
	w.Write([]byte((`{"result": {"status":"ok","msg":""}}\n`)))

}

func (h *Handler) processEvent(w http.ResponseWriter, r *http.Request, actionId sessionItem) {
	var jEvent JsonEvent

	// scanner := bufio.NewScanner(r.Body)
	// for scanner.Scan() {
	// 	ln := scanner.Text()
	// 	log.Println("scanner recieved", ln)
	// }

	err := json.NewDecoder(r.Body).Decode(&jEvent)
	if err != nil {
		log.Println(err)
	}

	switch jEvent.Method {
	case "event":
		// make sure the client id has ownership of the device id
		// and verify device exists

		// id = deviceid
		fmt.Println("!>>", jEvent.Data.Id, actionId.ClientId)
		// fmt.Println("!>>", actionId.ClientId)

		if !h.HomeManager.DeviceExistsWithClientId(jEvent.Data.Id, actionId.ClientId) {
			log.Println("Error: invalid device id")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error: incorrect device id"))
			return
		}
	}

	// TODO: validate the timestamp
	msg := event.EventMsg{
		Id:         jEvent.Data.Id,
		Timestamp:  jEvent.Data.Timestamp,
		Properties: jEvent.Data.Properties,
	}

	// add event to queue
	h.EventManager.AddEvent(msg)

}

func (h *Handler) callV1api(w http.ResponseWriter, r *http.Request, elements []string) {
	// r.RequestURI
	var id string
	var clientAuth login

	// TODO: this whole func is crap and needs a complete rewrite
	switch elements[2] {
	case "connect":
		// login
		err := json.NewDecoder(r.Body).Decode(&clientAuth)
		if err != nil {
			log.Printf("error reading body for %s: %s\n", r.URL.Path, err.Error())
		}

		newSession, ok := h.getSessionWithToken(clientAuth.Token)
		if !ok {
			return
		}
		fmt.Println("session id:", newSession)
		w.Header().Set("session", newSession)
		w.Write([]byte(``))

		return

	}

	sessionid := r.Header.Get("session")
	clientInfo, ok := h.lookupSessionClient(sessionid)
	if !ok {
		fmt.Println("bad session", clientInfo)
		return
	}

	if r.RequestURI == "/v1/register" {
		h.register(w, r)
		return
	}

	if len(elements) >= 4 {
		id = elements[3]
		// /v1/actions/uuid
		// fmt.Println(">>", id)
		val, _ := h.readActionID(id)
		// val, ok := h.deviceActionList[id]
		// if !ok {
		// 	log.Println("invalid id specified, id", id, "has not been registered")
		// }

		switch elements[2] {
		case "actions": // /v1/actions/uuid
			if id != clientInfo.actionId {
				log.Println("Error: invalid session")
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte("Error"))
				return
			}

			if val.inuse {
				val.done <- true
			}

			log.Println("client id", clientInfo.ClientId)
			val = waitActions{
				inuse: true,
				done:  make(chan bool),
				write: getWriter(w),
				resp:  &w,
				req:   r,
			}

			h.HomeManager.RegisterActionChannel(clientInfo.ClientId, &val)
			h.writeActionID(id, val)
			// h.deviceActionList[id] = val

			w.WriteHeader(http.StatusAccepted)
			val.write("")
			// this next line pauses
			h.waitClientActions(id, val)

		case "event": // /v1/event/uuid
			// if !val.inuse {
			// 	log.Println("Error: invalid event id")
			// 	w.WriteHeader(http.StatusBadRequest)
			// 	w.Write([]byte("Error: incorrect id"))
			// 	return
			// }
			h.processEvent(w, r, clientInfo)

		}
	}
}
