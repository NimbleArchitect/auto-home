package pluginManager

import "sync"

// var plugins map[string]Caller
type Plugin struct {
	lock sync.RWMutex
	list map[string]*PluginConnector
}

func (p *Plugin) Add(name string, connector *PluginConnector) {
	p.lock.Lock()
	if p.list == nil {
		p.list = make(map[string]*PluginConnector)
	}
	// fmt.Println("1>> add plugin", name)
	p.list[name] = connector
	p.lock.Unlock()
}

func (p *Plugin) All() map[string]*PluginConnector {
	p.lock.RLock()
	out := p.list
	p.lock.RUnlock()
	return out
}
