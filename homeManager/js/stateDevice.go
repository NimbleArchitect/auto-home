package js

type StateDevice struct {
	id         string
	name       string
	propSwitch map[string]jsSwitch
	propDial   map[string]jsDial
	propButton map[string]jsButton
	propText   map[string]jsText
}

// NewDevice initilises a new StateDevice object
func (r *JavascriptVM) NewDevice(id string, name string) StateDevice {
	dev := StateDevice{
		id:         id,
		name:       name,
		propSwitch: make(map[string]jsSwitch),
		propDial:   make(map[string]jsDial),
		propButton: make(map[string]jsButton),
		propText:   make(map[string]jsText),
	}

	return dev
}

// SaveDevice adds the StateDevice to the VMs deviceState list
// making it avaliable to the JS VM
func (r *JavascriptVM) SaveDevice(dev StateDevice) {

	r.deviceState[dev.id] = jsDevice{
		js:         r,
		Id:         dev.id,
		Name:       dev.name,
		propDial:   dev.propDial,
		propSwitch: dev.propSwitch,
		propButton: dev.propButton,
		propText:   dev.propText,
	}
}

// AddDial adds the dial property to the StateDevice
func (r *StateDevice) AddDial(id string, name string, value int, min int, max int) {
	r.propDial[id] = jsDial{
		Name:     name,
		Value:    value,
		min:      min,
		max:      max,
		previous: value,
	}
}

// AddSwitch adds the switch property to the StateDevice
func (r *StateDevice) AddSwitch(id string, name string, value bool, label string) {
	r.propSwitch[id] = jsSwitch{
		Name:     name,
		state:    value,
		Value:    label,
		previous: label,
	}
}

// AddButton adds the button property to the StateDevice
func (r *StateDevice) AddButton(id string, name string, value bool, label string) {
	r.propButton[id] = jsButton{
		Name:     name,
		Value:    value,
		label:    label,
		previous: label,
	}
}

// AddDial adds the dial property to the StateDevice
func (r *StateDevice) AddText(id string, name string, value string) {
	r.propText[id] = jsText{
		Name:     name,
		Value:    value,
		previous: value,
	}
}
