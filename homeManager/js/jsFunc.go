package js

import (
	"log"
	"net/rpc"
	"strings"

	"github.com/dop251/goja"
)

func (r *JavascriptVM) RunJS(deviceid string, fName string, props goja.Value) (goja.Value, error) {
	var jsHome jsHome
	var jsFunction goja.Value

	if r.runtime == nil {
		return nil, nil
	}

	parts := strings.Split(deviceid, "/")
	switch parts[0] {
	case "group":
		if val, ok := r.groupCode[deviceid]; ok {
			jsFunction = val.Get(fName)
		}
	default:
		if val, ok := r.deviceCode[deviceid]; ok {
			jsFunction = val.Get(fName)
		}
	}

	if jsFunction == nil {
		return nil, nil
	}

	call, ok := goja.AssertFunction(jsFunction)
	if !ok {
		// slient ignore as the function dosent exist in javascript
		log.Println("function", fName, "doesn't exist for device", deviceid, ", skipping")
		return nil, nil
	}

	jsHome.StopProcessing = FLAG_STOPPROCESSING
	jsHome.ContinueProcessing = FLAG_CONTINUEPROCESSING
	jsHome.GroupProcessing = FLAG_GROUPPROCESSING
	jsHome.devices = r.deviceState

	r.runtime.Set("home", jsHome)

	result, err := call(goja.Undefined(), props)
	if err != nil {
		log.Println(err)
	}

	return result, err
}

func (r *JavascriptVM) NewPlugin(client *rpc.Client, data map[string]interface{}) {
	// projectName := ""

	// for name, v := range data {
	// 	switch name {
	// 	case "name":
	// 		projectName = v.(string)
	// 	}
	// 	fmt.Println("name", name, v)
	// }

	// if len(projectName) > 0 {
	// 	globalPlugins[projectName] = *client
	// }
}
