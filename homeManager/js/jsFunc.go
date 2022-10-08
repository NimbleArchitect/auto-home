package js

import (
	"fmt"
	"log"
	"net/rpc"
	"strings"

	"github.com/dop251/goja"
)

func (r *JavascriptVM) RunJS(deviceid string, fName string, props goja.Value) (goja.Value, error) {
	var jsHome jsHome
	var jsFunction goja.Value

	var val *goja.Object
	var ok bool

	if r.runtime == nil {
		return nil, nil
	}

	parts := strings.Split(deviceid, "/")
	switch parts[0] {
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
	if !ok {
		// slient ignore as the function dosent exist in javascript
		log.Printf("function %s doesn't exist for %s, skipping", fName, deviceid)
		return nil, nil
	}

	jsHome.StopProcessing = FLAG_STOPPROCESSING
	jsHome.ContinueProcessing = FLAG_CONTINUEPROCESSING
	jsHome.GroupProcessing = FLAG_GROUPPROCESSING

	jsHome.devices = r.deviceState
	// fmt.Println("9>>", r.deviceState)

	jsHome.pluginList = r.pluginList

	r.runtime.Set("home", jsHome)

	result, err := call(goja.Undefined(), props)
	if err != nil {
		log.Println(err)
	}

	return result, err
}

func (r *JavascriptVM) NewPlugin(name string, vals *rpc.Client) {
	fmt.Println(">> add plugin", name, vals)
	r.pluginList[name] = vals

}
