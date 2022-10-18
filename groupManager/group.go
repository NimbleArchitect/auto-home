package groupManager

import (
	"fmt"
	"time"
)

type Group struct {
	Id          string
	Name        string
	Description string
	Devices     []string
	Groups      []string
	Users       []string

	CustomData   map[string]interface{}
	RepeatWindow int64

	repeatWindowTimeStamp time.Time
	repeatWindowDuration  time.Duration

	// ActionWriter func(s string) (int, error)
	// actionWriter ActionWriter
	// Groups      []*group

	// PropertySwitch map[string]*Switch
	// PropertyDial   map[string]*Dial
	// PropertyButton map[string]*Button
	// PropertyText   map[string]*Text

	// DialNames   []string
	// SwitchNames []string
	// ButtonNames []string
	// TextNames   []string

	// maxPropertyHistory int
	// Uploads []*Upload
}

func NewGroup() *Group {
	return &Group{}
}

func (d *Manager) Group(groupId string) (*Group, bool) {
	d.lock.RLock()
	out, ok := d.groups[groupId]
	d.lock.RUnlock()

	if ok {
		return out, true
	}
	return nil, false
}

func (m *Manager) SetGroup(groupId string, dev *Group) {
	m.lock.Lock()
	if _, ok := m.groups[groupId]; !ok {
		m.groupKeys = append(m.groupKeys, groupId)
	}

	m.groups[groupId] = dev
	m.lock.Unlock()
}

func (g *Group) Window(timestamp time.Time) bool {
	fmt.Println("1>>", g.RepeatWindow)
	fmt.Println("2>>", g.repeatWindowTimeStamp)
	fmt.Println("3>>", g.repeatWindowTimeStamp.Before(timestamp))
	fmt.Println("4>>", timestamp.Add(g.repeatWindowDuration))

	if g.RepeatWindow == 0 {
		return true
	}
	if g.repeatWindowTimeStamp.Before(timestamp) {
		g.repeatWindowTimeStamp = timestamp.Add(g.repeatWindowDuration)
		return true
	}
	// return true
	return false
}

func (g *Group) SetWindow(duration int64) {
	g.repeatWindowDuration = time.Duration(duration) * time.Millisecond
}
