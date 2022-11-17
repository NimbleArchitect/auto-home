package webHandle

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	event "server/eventManager"
	home "server/homeManager"
	"sync"
)

type Handler struct {
	ConfigPath string
	//deviceActionList map[string]*waitActions
	lockActionList sync.RWMutex

	//sessionTable     map[string]sessionItem
	//lockSessionTable sync.RWMutex

	//clientTable map[string]clientItem

	EventManager *event.Manager
	HomeManager  *home.Manager
	PublicPath   string
	FsHandle     http.Handler
	Address      string

	// TODO: add auth keys and variables
	userInfo map[string]clientItem

	//activeClients map[string]activeClientInfo // user map[clientID]userdetails

	// session maps session id to client id
	session map[string]sessionState // map[sessionid]clientid
}

// holds details needed for active connected clients
// type activeClientInfo struct {
// 	ActionID  string
// 	EventID   string
// 	Timestamp time.Time

// }

func (h *Handler) Shutdown() {
	h.lockActionList.Lock()
	for _, v := range h.session {
		if v.InUse {
			v.Done <- true
		}
	}

	h.lockActionList.Unlock()
}

func New(path string, publicPath string, evtMgr *event.Manager, homeMgr *home.Manager, hostAddress string) *Handler {

	return &Handler{
		ConfigPath:   path,
		EventManager: evtMgr,
		HomeManager:  homeMgr,
		FsHandle:     http.FileServer(http.Dir(publicPath)),
		Address:      hostAddress,
		// activeClients: make(map[string]activeClientInfo),
		session:  make(map[string]sessionState),
		userInfo: make(map[string]clientItem),
	}

}

func (h *Handler) SaveSystem() {
	log.Println("saving web configuration")

	// TODO: this has all been changed

	file, err := json.Marshal(h.userInfo)
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

	// TODO: this has all been changed

	h.userInfo["virtual.custom.light"] = clientItem{
		AuthKey: "secretclientid",
	}

	h.userInfo["virtual.custom.camera"] = clientItem{
		AuthKey: "randomcameradeviceuuid",
	}

	h.userInfo["device.front.door"] = clientItem{
		AuthKey: "randomfrontdoorid",
	}

	h.userInfo["com.ah.huehubv2"] = clientItem{
		AuthKey: "randomhuehubuuid",
	}

	// file, err := ioutil.ReadFile(path.Join(h.ConfigPath, "clients.json"))
	// if !errors.Is(err, os.ErrNotExist) {
	// 	if err != nil {
	// 		log.Panic("unable to read clients.json ", err)
	// 	}
	// 	err = json.Unmarshal(file, &h.userInfo)
	// 	if err != nil {
	// 		log.Panic("unable to read previous web state ", err)
	// 	}
	// }
}
