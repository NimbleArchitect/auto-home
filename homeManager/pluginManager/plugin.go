package pluginManager

import (
	log "server/logger"
	"sync"
)

// var plugins map[string]Caller
type Plugin struct {
	lock sync.RWMutex
	list map[string]*PluginConnector
}

func (p *Plugin) Add(name string, connector *PluginConnector) {
	log.Debug("lock")
	p.lock.Lock()
	if p.list == nil {
		p.list = make(map[string]*PluginConnector)
	}
	if con, ok := p.list[name]; ok {
		// if plugin exists we just update the connection and function list
		con.c = connector.c
		con.funcList = connector.funcList
		con.nextId = connector.nextId
	} else {
		// new item so we add the connector as is
		p.list[name] = connector
	}

	p.lock.Unlock()
	log.Debug("unlock")
}

func (p *Plugin) All() map[string]*PluginConnector {
	p.lock.RLock()
	out := p.list
	p.lock.RUnlock()
	return out
}

// Get returns the plugin named pluginName returns nil if the named plugin dosent exist
func (p *Plugin) Get(pluginName string) *PluginConnector {
	p.lock.RLock()
	out, ok := p.list[pluginName]
	p.lock.RUnlock()

	if ok {
		return out
	}
	return nil
}
