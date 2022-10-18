package groupManager

type iterator struct {
	groups map[string]*Group
	keys   []string
	index  int
	max    int
}

func (m *Manager) Iterate() *iterator {
	return &iterator{
		groups: m.groups,
		keys:   m.groupKeys,
		max:    len(m.groupKeys) - 1,
	}
}

func (i *iterator) Next() bool {
	if i.index < i.max {
		i.index++
		return true
	}
	return false
}

func (i *iterator) Get() *Group {
	return i.groups[i.keys[i.index]]
}
