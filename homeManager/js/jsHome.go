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

func (d *jsHome) GetDeviceByName(name string) jsDevice {
	for _, v := range d.devices {
		if v.Name == name {
			return v
		}
	}

	return jsDevice{}
}

func (d *jsHome) GetDeviceById(deviceId string) jsDevice {
	for _, v := range d.devices {
		if v.Id == deviceId {
			return v
		}
	}

	return jsDevice{}
}

func (d *jsHome) GetDevices() []jsDevice {
	var out []jsDevice
	for _, v := range d.devices {
		out = append(out, v)
	}

	return out
}

func (d *jsHome) GetDevicesStartName(s string) []jsDevice {
	var out []jsDevice

	for k, v := range d.devices {
		if strings.HasPrefix(k, s) {
			out = append(out, v)
		}
	}

	return out
}

func (d *jsHome) GetGroupByName(s string) []jsDevice {
	var out []jsDevice

	for k, v := range d.devices {
		if strings.HasPrefix(k, s) {
			out = append(out, v)
		}
	}

	return out
}

// Sleep pauses the vm execution for the specified number of seconds
func (d *jsHome) Sleep(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}

// Countdown creates a counter thats identified by name and calls function after mSec milliseconds
//
//	if called multiple times the counter is reset, if mSec is 0 the counter is removed
func (d *jsHome) Countdown(name string, mSec int, function goja.Value) {
	var jsCall goja.Callable
	var ok bool

	if !goja.IsUndefined(function) {
		jsCall, ok = goja.AssertFunction(function)
	}

	d.vm.waitGroup.TryLock()

	d.vm.global.SetTimer(name, mSec, func(success bool) {
		if !success {
			d.vm.waitGroup.Unlock()
			return
		}

		// we skip the js call if function isnt defined or was invalid
		if ok && !goja.IsUndefined(function) {
			jsCall(goja.Undefined())
		}

		d.vm.waitGroup.Unlock()
	})
}
