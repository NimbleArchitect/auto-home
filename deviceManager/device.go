package deviceManager

type Device struct {
	Id          string
	Name        string
	Description string
	Help        string
	ClientId    string
	Groups      []*group

	PropertySwitch map[string]*Switch
	PropertyDial   map[string]*Dial
	PropertyButton map[string]*Button
	PropertyText   map[string]*Text

	DialNames   []string
	SwitchNames []string
	ButtonNames []string
	TextNames   []string

	Uploads []*Upload
}

func NewDevice() *Device {
	return &Device{
		PropertyDial:   make(map[string]*Dial),
		PropertySwitch: make(map[string]*Switch),
		PropertyButton: make(map[string]*Button),
		PropertyText:   make(map[string]*Text),
	}
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

func (d *Manager) SetDevice(deviceId string, dev *Device) {
	d.lock.Lock()
	if _, ok := d.devices[deviceId]; !ok {
		d.deviceKeys = append(d.deviceKeys, deviceId)
	}
	d.devices[deviceId] = dev

	d.lock.Unlock()
}
