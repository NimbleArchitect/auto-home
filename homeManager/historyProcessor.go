package home

import "sync"

type historyProcessor struct {
	lock    sync.RWMutex
	history []eventHistory
	max     int
}

func (h *historyProcessor) Add(event eventHistory) {
	h.lock.Lock()
	h.history = append(h.history, event)

	if len(h.history) > h.max {
		s := len(h.history) - h.max
		h.history = h.history[s:h.max]
	}
	h.lock.Unlock()
}

func (h *historyProcessor) Latest() eventHistory {
	h.lock.RLock()
	out := h.history[len(h.history)-1]
	h.lock.RUnlock()
	return out
}
