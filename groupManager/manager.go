package groupManager

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
)

type Manager struct {
	lock *sync.RWMutex
	// lock    *lock
	groups map[string]*Group

	// window map[string]*duration

	// maxPropertyHistory int
	groupKeys []string
}

type onDiskGroup struct {
	Id           string
	Name         string
	Description  string
	Devices      []string
	Groups       []string
	Users        []string
	RepeatWindow int64
}

func New() *Manager {
	return &Manager{
		lock:   &sync.RWMutex{},
		groups: make(map[string]*Group),
		// maxPropertyHistory: maxPropertyHistory,
	}
}

func (m *Manager) Save() {
	log.Println("saving groups")

	groupList := make(map[string]onDiskGroup)

	for key, g := range m.groups {
		groupList[key] = onDiskGroup{
			Id:           g.Id,
			Name:         g.Name,
			Description:  g.Description,
			Devices:      g.Devices,
			Groups:       g.Groups,
			Users:        g.Users,
			RepeatWindow: g.RepeatWindow,
		}
	}
	file, err := json.Marshal(groupList)
	if err != nil {
		log.Println("unable to serialize groups", err)
	}
	err = os.WriteFile("groups.json", file, 0640)
	if err != nil {
		log.Println("unable to write groups.json", err)
	}

}

func (m *Manager) Load() {
	var groupList map[string]onDiskGroup

	file, err := os.ReadFile("groups.json")
	if !errors.Is(err, os.ErrNotExist) {
		if err != nil {
			log.Panic("unable to read groups.json ", err)
		}
		err = json.Unmarshal(file, &groupList)
		if err != nil {
			log.Panic("unable to read previous group state ", err)
		}
	}

	for id, group := range groupList {
		g := NewGroup()
		g.Id = group.Id
		g.Name = group.Name
		g.Description = group.Description
		g.Devices = group.Devices
		g.Groups = group.Groups
		g.Users = group.Users
		g.RepeatWindow = group.RepeatWindow
		g.SetWindow(group.RepeatWindow)
		m.SetGroup(id, g)
	}
}
