package js

import (
	"server/booltype"
	"server/deviceManager"
	"server/logger"
	"strconv"

	"github.com/dop251/goja"
)

type jsDevice struct {
	js         *JavascriptVM
	Id         string
	Name       string
	propDial   map[string]jsDial
	propSwitch map[string]jsSwitch
	propButton map[string]jsButton
	propText   map[string]jsText
	liveDevice *deviceManager.Device // points back to the device stored within the deviceManager
}

func (d *jsDevice) GetDial(name string) interface{} {
	if val, ok := d.propDial[name]; ok {
		// create a js object so we cn add a new property
		jsVal := d.js.runtime.ToValue(val)
		jsObj := d.js.runtime.CreateObject(jsVal.ToObject(d.js.runtime))
		// add a readonly .latest propery that gets the live device property value
		if d.liveDevice != nil {
			jsObj.DefineAccessorProperty("latest", d.js.runtime.ToValue(func() interface{} {
				if val, ok := d.liveDevice.DialValue(name); ok {
					return val
				}
				return nil
			}),
				nil, goja.FLAG_FALSE, goja.FLAG_FALSE)
		}

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
		if d.liveDevice != nil {
			jsObj.DefineAccessorProperty("latest", d.js.runtime.ToValue(func() interface{} {
				if val, ok := d.liveDevice.SwitchValue(name); ok {
					return jsBool{
						s: val.String(),
						b: val.GetBool(),
					}
				}
				return nil
			}),
				nil, goja.FLAG_FALSE, goja.FLAG_FALSE)
		}

		// also add previous as a readonly property
		jsObj.DefineAccessorProperty("previous", d.js.runtime.ToValue(func() jsBool {
			boolVal := booltype.BoolType{}
			boolVal.Set(val.previous)
			return jsBool{
				s: boolVal.String(),
				b: boolVal.GetBool(),
			}
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
		if d.liveDevice != nil {
			jsObj.DefineAccessorProperty("latest", d.js.runtime.ToValue(func() interface{} {
				if val, ok := d.liveDevice.ButtonValue(name); ok {
					return jsBool{
						s: val.String(),
						b: val.GetBool(),
					}
				}
				return nil
			}),
				nil, goja.FLAG_FALSE, goja.FLAG_FALSE)
		}

		// also add previous as a readonly property
		jsObj.DefineAccessorProperty("previous", d.js.runtime.ToValue(func() jsBool {
			boolVal := booltype.BoolType{}
			boolVal.Set(val.previous)
			return jsBool{
				s: boolVal.String(),
				b: boolVal.GetBool(),
			}
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
		if d.liveDevice != nil {
			jsObj.DefineAccessorProperty("latest", d.js.runtime.ToValue(func() interface{} {
				if val, ok := d.liveDevice.TextValue(name); ok {
					return val
				}
				return nil
			}),
				nil, goja.FLAG_FALSE, goja.FLAG_FALSE)
		}

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
	log := logger.New("Set", &debugLevel)

	for _, v := range d.propDial {
		if v.Name == name {
			i, err := strconv.Atoi(value)
			if err != nil {
				log.Error("Not a valid number")
				continue
			}
			d.liveDevice.SetDialValue(name, i)
			d.liveDevice.WriteDialValue(name, i)
		}
	}

	for _, v := range d.propSwitch {
		if v.Name == name {
			d.liveDevice.SetSwitchValue(name, value)
			d.liveDevice.WriteSwitchValue(name, value)
		}
	}

	for _, v := range d.propButton {
		if v.Name == name {
			d.liveDevice.SetButtonValue(name, value)
			d.liveDevice.WriteButtonValue(name, value)
		}
	}

	for _, v := range d.propText {
		if v.Name == name {
			d.liveDevice.SetTextValue(name, value)
			d.liveDevice.WriteTextValue(name, value)
		}
	}
}
