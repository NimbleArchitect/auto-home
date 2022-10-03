package js

import (
	"log"
	"strconv"

	"github.com/dop251/goja"
)

type jsDevice struct {
	js         *JavascriptVM
	Id         string
	Name       string
	groups     map[string]jsGroup
	propDial   map[string]jsDial
	propSwitch map[string]jsSwitch
	propButton map[string]jsButton
	propText   map[string]jsText
}

func (d *jsDevice) GetDial(name string) interface{} {
	if val, ok := d.propDial[name]; ok {
		// create a js object so we cn add a new property
		jsVal := d.js.runtime.ToValue(val)
		jsObj := d.js.runtime.CreateObject(jsVal.ToObject(d.js.runtime))
		// add a readonly .latest propery that gets the live device property value
		jsObj.DefineAccessorProperty("latest", d.js.runtime.ToValue(func() interface{} {
			if val, ok := d.js.Updater.GetDialValue(d.Id, name); ok {
				return val
			}
			return nil
		}),
			nil, goja.FLAG_FALSE, goja.FLAG_FALSE)

		// we also add previous as a readonly property
		jsObj.DefineAccessorProperty("previous", d.js.runtime.ToValue(func() int {
			return val.previous
		}),
			nil, goja.FLAG_FALSE, goja.FLAG_FALSE)

		return jsObj
	}

	return nil
}

func (d *jsDevice) GetSwitch(name string) interface{} {
	if val, ok := d.propSwitch[name]; ok {
		// create a js object so we cn add a new property
		jsVal := d.js.runtime.ToValue(val)
		jsObj := d.js.runtime.CreateObject(jsVal.ToObject(d.js.runtime))
		// add a readonly .latest propery that gets the live device property value
		jsObj.DefineAccessorProperty("latest", d.js.runtime.ToValue(func() interface{} {
			if val, ok := d.js.Updater.GetSwitchValue(d.Id, name); ok {
				return val
			}
			return nil
		}),
			nil, goja.FLAG_FALSE, goja.FLAG_FALSE)

		// also add previous as a readonly property
		jsObj.DefineAccessorProperty("previous", d.js.runtime.ToValue(func() string {
			return val.previous
		}),
			nil, goja.FLAG_FALSE, goja.FLAG_FALSE)

		return jsObj
	}

	return nil
}

func (d *jsDevice) GetButton(name string) interface{} {
	if val, ok := d.propButton[name]; ok {
		// create a js object so we cn add a new property
		jsVal := d.js.runtime.ToValue(val)
		jsObj := d.js.runtime.CreateObject(jsVal.ToObject(d.js.runtime))
		// add a readonly .latest propery that gets the live property value from the device
		jsObj.DefineAccessorProperty("latest", d.js.runtime.ToValue(func() interface{} {
			if val, ok := d.js.Updater.GetButtonValue(d.Id, name); ok {
				return val
			}
			return nil
		}),
			nil, goja.FLAG_FALSE, goja.FLAG_FALSE)

		// also add previous as a readonly property
		jsObj.DefineAccessorProperty("previous", d.js.runtime.ToValue(func() string {
			return val.previous
		}),
			nil, goja.FLAG_FALSE, goja.FLAG_FALSE)

		return jsObj
	}

	return nil
}

func (d *jsDevice) GetText(name string) interface{} {
	if val, ok := d.propText[name]; ok {
		// create a js object so we cn add a new property
		jsVal := d.js.runtime.ToValue(val)
		jsObj := d.js.runtime.CreateObject(jsVal.ToObject(d.js.runtime))
		// add a readonly .latest propery that gets the live device property value
		jsObj.DefineAccessorProperty("latest", d.js.runtime.ToValue(func() interface{} {
			if val, ok := d.js.Updater.GetTextValue(d.Id, name); ok {
				return val
			}
			return nil
		}),
			nil, goja.FLAG_FALSE, goja.FLAG_FALSE)

		// also add previous as a readonly property
		jsObj.DefineAccessorProperty("previous", d.js.runtime.ToValue(func() string {
			return val.previous
		}),
			nil, goja.FLAG_FALSE, goja.FLAG_FALSE)

		return jsObj
	}

	return nil
}

func (d *jsDevice) Get(name string) interface{} {

	if obj := d.GetDial(name); obj != nil {
		return obj
	}

	if obj := d.GetSwitch(name); obj != nil {
		return obj
	}

	if obj := d.GetButton(name); obj != nil {
		return obj
	}

	if obj := d.GetText(name); obj != nil {
		return obj
	}
	return nil
}

func (d *jsDevice) Set(name string, value string) {

	for _, v := range d.propDial {
		if v.Name == name {
			i, err := strconv.Atoi(value)
			if err != nil {
				log.Println("Not a valid number")
				continue
			}
			d.js.Updater.UpdateDial(d.Id, name, i)
		}
	}

	for _, v := range d.propSwitch {
		if v.Name == name {
			d.js.Updater.UpdateSwitch(d.Id, name, value)
		}
	}

	for _, v := range d.propButton {
		if v.Name == name {
			d.js.Updater.UpdateButton(d.Id, name, value)
		}
	}

	for _, v := range d.propText {
		if v.Name == name {
			d.js.Updater.UpdateText(d.Id, name, value)
		}
	}
}
