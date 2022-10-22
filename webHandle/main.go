package webHandle

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	event "server/eventManager"
	home "server/homeManager"
	"sync"
)

type Handler struct {
	ConfigPath       string
	deviceActionList map[string]waitActions
	lockActionList   sync.RWMutex

	sessionTable     map[string]sessionItem
	lockSessionTable sync.RWMutex

	clientTable map[string]clientItem

	EventManager *event.Manager
	HomeManager  *home.Manager
	PublicPath   string
	FsHandle     http.Handler
	Address      string
}

func (h *Handler) Shutdown() {
	h.lockActionList.Lock()
	for _, v := range h.deviceActionList {
		if v.inuse {
			v.done <- true
		}
	}

	h.lockActionList.Unlock()
}

func (h *Handler) unRegisterClientAction(id string) {

	if val, ok := h.readActionID(id); ok {
		val.done <- true
		h.deleteActionID(id)
		// delete(h.deviceActionList, id)
	}
	log.Println("wait action", id, "closed")
}

func (h *Handler) waitClientActions(val string, wait waitActions) {

	ctx := wait.req.Context()

	defer h.unRegisterClientAction(val)

	select {
	case <-ctx.Done():
		log.Println("finished waitClientActions")
		http.Error(*wait.resp, ctx.Err().Error(), http.StatusInternalServerError)
		tmp, _ := h.readActionID(val)
		// tmp := h.deviceActionList[val]
		tmp.inuse = false
		h.writeActionID(val, tmp)
		// h.deviceActionList[val] = tmp
	case <-wait.done:
		log.Println("waitClientActions wait.done")
		tmp, _ := h.readActionID(val)
		// tmp := h.deviceActionList[val]
		tmp.inuse = false
		h.writeActionID(val, tmp)
		// h.deviceActionList[val] = tmp
	}

}

func (h *Handler) SaveSystem() {
	log.Println("saving web configuration")

	file, err := json.Marshal(h.clientTable)
	if err != nil {
		log.Println("unable to serialize clients", err)
	}
	err = ioutil.WriteFile(path.Join(h.ConfigPath, "clients.json"), file, 0640)
	if err != nil {
		log.Println("unable to write clients.json", err)
	}
}

func (h *Handler) LoadSystem() {
	log.Println("loading web configuration")

	file, err := ioutil.ReadFile(path.Join(h.ConfigPath, "clients.json"))
	if !errors.Is(err, os.ErrNotExist) {
		if err != nil {
			log.Panic("unable to read clients.json ", err)
		}
		err = json.Unmarshal(file, &h.clientTable)
		if err != nil {
			log.Panic("unable to read previous web state ", err)
		}
	}
}
