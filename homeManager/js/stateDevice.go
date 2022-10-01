package js

type StateDevice struct {
	id         string
	name       string
	propSwitch map[string]jsSwitch
	propDial   map[string]jsDial
}

// NewDevice initilises a new StateDevice object
func (r *JavascriptVM) NewDevice(id string, name string) StateDevice {
	dev := StateDevice{
		id:         id,
		name:       name,
		propSwitch: make(map[string]jsSwitch),
		propDial:   make(map[string]jsDial),
	}

	return dev
}

// SaveDevice adds the StateDevice to the VMs deviceState list
// making it avaliable to the JS VM
func (r *JavascriptVM) SaveDevice(dev StateDevice) {
	if len(r.deviceState) == 0 {
		r.deviceState = make(map[string]jsDevice)
	}

	r.deviceState[dev.id] = jsDevice{
		vm:         r,
		Id:         dev.id,
		Name:       dev.name,
		propDial:   dev.propDial,
		propSwitch: dev.propSwitch,
	}
}

// AddDial adds the dial property to the StateDevice
func (r *StateDevice) AddDial(id string, name string, value int) {
	r.propDial[id] = jsDial{
		Name:  name,
		Value: value,
	}
}

// AddSwitch adds the switch property to the StateDevice
func (r *StateDevice) AddSwitch(id string, name string, value bool, label string) {
	r.propSwitch[id] = jsSwitch{
		Name:  name,
		Value: value,
		label: label,
	}
}
