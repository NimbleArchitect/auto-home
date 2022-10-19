package deviceManager

type iterator struct {
	devices map[string]*Device
	keys    []string
	index   int
	max     int
}

func (m *Manager) Iterate() *iterator {
	return &iterator{
		devices: m.devices,
		keys:    m.deviceKeys,
		max:     len(m.deviceKeys) - 1,
		index:   -1,
	}
}

func (i *iterator) Next() bool {
	if i.index < i.max {
		i.index++
		return true
	}
	return false
}

func (i *iterator) Get() (string, *Device) {
	name := i.keys[i.index]
	return name, i.devices[name]
}
