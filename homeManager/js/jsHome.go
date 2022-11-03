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

	countdown goja.Object
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

func (d *jsHome) Sleep(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}

// func (d *jsHome) Countdown(name string, mSec int, function goja.Value) {
// 	jscall, _ := goja.AssertFunction(function)

// 	d.vm.waitGroup.TryLock()

// 	d.vm.global.SetTimer(name, mSec, func() {
// 		jscall(goja.Undefined())

// 		d.vm.waitGroup.Unlock()
// 	})
// }
