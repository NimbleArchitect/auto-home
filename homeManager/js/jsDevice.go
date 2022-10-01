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
}

func (d *jsDevice) Get(name string) interface{} {

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

	if val, ok := d.propSwitch[name]; ok {
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

		// also add previous as a readonly property
		jsObj.DefineAccessorProperty("previous", d.js.runtime.ToValue(func() string {
			return val.previous
		}),
			nil, goja.FLAG_FALSE, goja.FLAG_FALSE)

		return jsObj
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

}
