package groupManager

import (
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
	if g.RepeatWindow == 0 {
		return false
	}
	if g.repeatWindowTimeStamp.Before(timestamp) {
		g.repeatWindowTimeStamp = timestamp.Add(g.repeatWindowDuration)
		return true
	}
	return false
}

func (g *Group) SetWindow(duration int64) {
	g.repeatWindowDuration = time.Duration(duration) * time.Millisecond
}
