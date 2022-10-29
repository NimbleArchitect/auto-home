package js

import (
	"log"
	"server/deviceManager"
	"server/homeManager/pluginManager"
	"strings"
	"time"
	"unicode"

	"github.com/dop251/goja"
)

type JavascriptVM struct {
	runtime     *goja.Runtime
	deviceCode  map[string]*goja.Object
	deviceState map[string]jsDevice
	groupCode   map[string]*goja.Object
	groups      map[string]jsGroup
	userCode    map[string]*goja.Object
	pluginCode  map[string]*goja.Object
	pluginList  *pluginManager.Plugin // plugin connections, shared across all VMs
	// users      map[string]jsUser
	Updater DeviceUpdator
}

// RunJS loads the js object attached to the specified deviceId and runs the function fName passing in props as an argument
func (r *JavascriptVM) RunJS(deviceid string, fName string, props goja.Value) (goja.Value, error) {
	var jsHome jsHome
	var jsFunction goja.Value

	var err error
	var result goja.Value

	var val *goja.Object
	var ok bool

	if r.runtime == nil {
		return nil, nil
	}

	parts := strings.Split(deviceid, "/")
	switch parts[0] {
	case "plugin":
		val, ok = r.pluginCode[deviceid]
	case "group":
		val, ok = r.groupCode[deviceid]
	case "user":
		val, ok = r.userCode[deviceid]
	case "device":
		fallthrough
	default:
		val, ok = r.deviceCode[deviceid]
	}

	if ok {
		jsFunction = val.Get(fName)
		if jsFunction == nil {
			return nil, nil
		}
	}

	call, ok := goja.AssertFunction(jsFunction)
	if !ok || call == nil {
		// slient ignore as the function dosent exist in javascript
		log.Printf("function %s doesn't exist for %s, skipping", fName, deviceid)
		return nil, nil
	}

	jsHome.StopProcessing = FLAG_STOPPROCESSING
	jsHome.ContinueProcessing = FLAG_CONTINUEPROCESSING
	jsHome.GroupProcessing = FLAG_GROUPPROCESSING
	jsHome.PreventUpdate = FLAG_PREVENTUPDATE

	jsHome.devices = r.deviceState

	r.runtime.Set("home", jsHome)

	if props == nil {
		log.Printf("calling %s/%s with no arguments\n", deviceid, fName)
		result, err = call(goja.Undefined())
	} else {
		log.Printf("calling %s/%s with one argument\n", deviceid, fName)
		result, err = call(goja.Undefined(), props)
	}
	if err != nil {
		log.Println("err", err)
	}

	return result, err
}

func (r *JavascriptVM) RunJSGroup(groupId string, props JSPropsList) (int64, error) {
	val, err := r.RunJS("group/"+groupId, "onchange", r.runtime.ToValue(props))
	if err != nil {
		return 0, err
	} else {
		if val == nil {
			return 0, nil
		}
		return r.runtime.ToValue(val).ToInteger(), nil
	}

}

func (r *JavascriptVM) RunJSPlugin(pluginName string, fName string, args map[string]interface{}) (int64, error) {
	obj := r.runtime.NewObject()
	for rawName, value := range args {
		name := []rune(rawName)
		name[0] = unicode.ToLower(name[0])

		obj.Set(string(name), value)
	}

	val, err := r.RunJS("plugin/"+pluginName, fName, obj)
	if err != nil {
		return 0, err
	} else {
		if val == nil {
			return 0, nil
		}
		return r.runtime.ToValue(val).ToInteger(), nil
	}

}

// Process main entry point after a trigger, this allows processin gthe event data
func (r *JavascriptVM) Process(deviceid string, timestamp time.Time, props JSPropsList) {
	var dev jsDevice

	log.Println("process triggered")

	dev.propSwitch = make(map[string]jsSwitch)
	dev.propDial = make(map[string]jsDial)
	dev.propButton = make(map[string]jsButton)
	dev.propText = make(map[string]jsText)

	// log.Println("state:", m.devices)

	// lookup changes and trigger change notifications
	r.processOnTrigger(deviceid, timestamp, props, &dev)

	// TODO: not sure this is the correct order as it depends on if we wnat groups to return a no further processing argument
	continueFlag := r.processGroupChange(deviceid, props)
	r.processOnChange(deviceid, &dev, continueFlag)
}

// SaveDeviceState copies the current device sate and properties to the vm ready for processing/usage
func (r *JavascriptVM) SaveDeviceState(devices *deviceManager.Manager) {

	iterator := devices.Iterate()

	for iterator.Next() {
		deviceId, dev := iterator.Get()
		newDev := jsDevice{
			js:         r,
			Id:         dev.Id,
			Name:       dev.Name,
			propSwitch: make(map[string]jsSwitch),
			propDial:   make(map[string]jsDial),
			propButton: make(map[string]jsButton),
			propText:   make(map[string]jsText),
			liveDevice: dev,
		}
		for key, property := range dev.DialAsMap() {
			newDev.propDial[key] = jsDial{
				Name:     property.Name,
				Value:    property.Value,
				min:      property.Min,
				max:      property.Max,
				previous: property.Value,
			}
		}

		for key, property := range dev.SwitchAsMap() {
			newDev.propSwitch[key] = jsSwitch{
				Name:     property.Name,
				Value:    property.Value.String(),
				state:    property.Value.GetBool(),
				previous: property.Value.String(),
			}
		}

		for key, property := range dev.ButtonAsMap() {
			newDev.propButton[key] = jsButton{
				Name:     property.Name,
				Value:    property.Value.GetBool(),
				label:    property.Value.String(),
				previous: property.Value.String(),
			}
		}

		for key, property := range dev.TextAsMap() {
			newDev.propText[key] = jsText{
				Name:     property.Name,
				Value:    property.Value,
				previous: property.Value,
			}
		}
		r.deviceState[deviceId] = newDev
	}

}
