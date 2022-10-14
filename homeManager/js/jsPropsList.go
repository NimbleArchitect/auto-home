package js

type JSPropsList struct {
	propDial   map[string]jsDial
	propSwitch map[string]jsSwitch
	propButton map[string]jsButton
	propText   map[string]jsText
}

func NewJSDevice() JSPropsList {
	return JSPropsList{
		propSwitch: make(map[string]jsSwitch),
		propDial:   make(map[string]jsDial),
		propButton: make(map[string]jsButton),
		propText:   make(map[string]jsText),
	}
}

func (d *JSPropsList) AddDial(name string, prop jsDial) {
	d.propDial[name] = prop
}

func (d *JSPropsList) AddSwitch(name string, prop jsSwitch) {
	d.propSwitch[name] = prop
}

func (d *JSPropsList) AddButton(name string, prop jsButton) {
	d.propButton[name] = prop
}

func (d *JSPropsList) AddText(name string, prop jsText) {
	d.propText[name] = prop
}
