package js

import (
	"strings"
	"time"

	"github.com/dop251/goja"
)

type jsHome struct {
	vm      *JavascriptVM
	devices map[string]jsDevice
	// groups             map[string]jsGroup
	// pluginList         map[string]*pluginManager.Caller
	StopProcessing     int
	GroupProcessing    int
	ContinueProcessing int
	PreventUpdate      int
}

func (h *jsHome) GetDeviceByName(name string) jsDevice {
	for _, v := range h.devices {
		if v.Name == name {
			return v
		}
	}

	return jsDevice{}
}

func (h *jsHome) GetDeviceById(deviceId string) jsDevice {
	for _, v := range h.devices {
		if v.Id == deviceId {
			return v
		}
	}

	return jsDevice{}
}

func (h *jsHome) GetDevices() []jsDevice {
	var out []jsDevice
	for _, v := range h.devices {
		out = append(out, v)
	}

	return out
}

func (h *jsHome) GetDevicesStartName(s string) []jsDevice {
	var out []jsDevice

	for _, v := range h.devices {
		if strings.HasPrefix(v.Name, s) {
			out = append(out, v)
		}
	}

	return out
}

// func (h *jsHome) GetGroupByName(s string) []jsDevice {
// 	var out []jsDevice

// 	for k, v := range h.devices {
// 		if strings.HasPrefix(k, s) {
// 			out = append(out, v)
// 		}
// 	}

// 	return out
// }

// Sleep pauses the vm execution for the specified number of seconds
func (h *jsHome) Sleep(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}

// Countdown creates a counter thats identified by name and calls function after sec seconds
//
// if called multiple times the counter is reset, if sec is 0 the counter is removed,
// sec is a float and supports fractional seconds
func (h *jsHome) Countdown(name string, sec float64, function goja.Value) {
	var jsCall goja.Callable
	var ok bool

	if !goja.IsUndefined(function) {
		jsCall, ok = goja.AssertFunction(function)
	}

	h.vm.waitLock.TryLock()

	h.vm.global.SetTimer(name, sec, func(success bool) {
		if !success {
			h.vm.waitLock.Unlock()
			return
		}

		// we skip the js call if function isnt defined or was invalid
		if ok && !goja.IsUndefined(function) {
			jsCall(goja.Undefined())
		}

		h.vm.waitLock.Unlock()
	})
}

func (h *jsHome) SetGlobal(name string, value interface{}) {
	h.vm.global.SetVariable(name, value)
}

func (h *jsHome) GetGlobal(name string) interface{} {
	return h.vm.global.GetVariable(name)
}
