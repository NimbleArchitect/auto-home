package js

import (
	"server/deviceManager"
	"server/globals"
	"server/homeManager/pluginManager"
	log "server/logger"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/dop251/goja"
)

type JavascriptVM struct {
	waitGroup   sync.Mutex
	global      *globals.Global
	runtime     *goja.Runtime
	deviceCode  map[string]*goja.Object // list of compiled javascript device code that has been registered using the javascript set function, used to store onchange functions
	deviceState map[string]jsDevice     //list of devices that the javascrip VM can use
	groupCode   map[string]*goja.Object // list of compiled javascript group code that has been registered using the javascript set function
	groups      map[string]jsGroup
	userCode    map[string]*goja.Object
	plugins     map[string]*goja.Object //vm's copy of all plugins attached to the plugin object
	pluginCode  map[string]*goja.Object // list of compiled javascript plugin code that has been registered using the javascript set function
	pluginList  *pluginManager.Plugin   // plugin connections, shared across all VMs
	// users      map[string]jsUser
	Updater DeviceUpdator
}

func (r *JavascriptVM) Wait() {
	r.waitGroup.Lock()

	r.waitGroup.Unlock()
}

// RunJS loads the js object attached to the specified deviceId and runs the function fName passing in props as an argument
func (r *JavascriptVM) RunJS(deviceid string, fName string, props goja.Value) (goja.Value, error) {
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
		log.Infof("function %s doesn't exist for %s, skipping\n", fName, deviceid)
		return nil, nil
	}

	r.setJsGlobal()

	if props == nil {
		log.Infof("calling %s/%s with no arguments\n", deviceid, fName)
		result, err = call(goja.Undefined())
	} else {
		log.Infof("calling %s/%s with one argument\n", deviceid, fName)
		result, err = call(goja.Undefined(), props)
	}
	if err != nil {
		log.Error(err)
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

// Process main entry point after a trigger, this allows processing the event data
func (r *JavascriptVM) Process(deviceid string, timestamp time.Time, props JSPropsList) {
	var dev jsDevice

	log.Info("process triggered")

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
				Value:    property.Value.String(),
				state:    property.Value.GetBool(),
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
