package js

import (
	"log"
	"net/rpc"
	"time"

	"github.com/dop251/goja"
)

type JavascriptVM struct {
	runtime     *goja.Runtime
	deviceCode  map[string]*goja.Object
	deviceState map[string]jsDevice
	groupCode   map[string]*goja.Object
	groups      map[string]jsGroup
	userCode    map[string]*goja.Object
	pluginList  map[string]*rpc.Client
	// users      map[string]jsUser
	Updater DeviceUpdator
}

func (r *JavascriptVM) RunJSGroupAction(groupId string, fnName string, props []map[string]interface{}) (interface{}, error) {
	// var dev jsDevice

	log.Println("group action triggered:", groupId, fnName)

	// dev.propSwitch = make(map[string]jsSwitch)
	// dev.propDial = make(map[string]jsDial)
	// dev.propButton = make(map[string]jsButton)
	// dev.propText = make(map[string]jsText)

	// log.Println("state:", m.devices)

	// lookup changes and trigger change notifications
	out, err := r.RunJS("group/"+groupId, fnName, r.runtime.ToValue(props))
	if err != nil {
		log.Println(err)
	}

	return out, nil
}

// Process main entry point after a trigger, this allows processin gthe event data
func (r *JavascriptVM) Process(deviceid string, timestamp time.Time, props []map[string]interface{}) {
	var dev jsDevice

	log.Println("event triggered")

	dev.propSwitch = make(map[string]jsSwitch)
	dev.propDial = make(map[string]jsDial)
	dev.propButton = make(map[string]jsButton)
	dev.propText = make(map[string]jsText)

	// log.Println("state:", m.devices)

	// lookup changes and trigger change notifications
	r.processOnTrigger(deviceid, timestamp, props, &dev)

	// TODO: not sure this is the correct order as it depends on if we wnat groups to return a no further processing argument
	continueFlag := r.processGroupChange(deviceid, timestamp, props, &dev)
	if continueFlag != FLAG_STOPPROCESSING {
		r.processOnChange(deviceid, timestamp, props, &dev)
	}

}
