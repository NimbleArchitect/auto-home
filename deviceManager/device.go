package deviceManager

type Device struct {
	Id          string
	Name        string
	Description string
	Help        string
	ClientId    string
	// Groups      []*group

	PropertySwitch map[string]*Switch
	PropertyDial   map[string]*Dial
	PropertyButton map[string]*Button
	PropertyText   map[string]*Text

	DialNames   []string
	SwitchNames []string
	ButtonNames []string
	TextNames   []string

	repeatWindow       map[string]int64
	maxPropertyHistory int
	// Uploads []*Upload
}

func NewDevice(maxPropertyHistory int) *Device {
	return &Device{
		PropertyDial:   make(map[string]*Dial),
		PropertySwitch: make(map[string]*Switch),
		PropertyButton: make(map[string]*Button),
		PropertyText:   make(map[string]*Text),

		maxPropertyHistory: maxPropertyHistory,
	}
}

func (d *Device) MaxHistory() int {
	return d.maxPropertyHistory
}

func (d *Manager) Device(deviceId string) (*Device, bool) {
	d.lock.RLock()
	out, ok := d.devices[deviceId]
	d.lock.RUnlock()

	if ok {
		return out, true
	}
	return nil, false
}

func (d *Manager) HasDevice(deviceId string) bool {
	d.lock.RLock()
	_, ok := d.devices[deviceId]
	d.lock.RUnlock()

	return ok
}

func (m *Manager) SetDevice(deviceId string, dev *Device) {
	m.lock.Lock()
	if _, ok := m.devices[deviceId]; !ok {
		m.deviceKeys = append(m.deviceKeys, deviceId)
	}

	m.devices[deviceId] = dev
	m.lock.Unlock()
}
