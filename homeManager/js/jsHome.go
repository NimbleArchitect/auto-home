package js

import (
	"fmt"
	"net/rpc"
	"strings"
	"time"
)

const (
	FLAG_CONTINUEPROCESSING = iota
	FLAG_GROUPPROCESSING
	FLAG_STOPPROCESSING
)

type jsHome struct {
	devices map[string]jsDevice
	// groups             map[string]jsGroup
	pluginList         map[string]*rpc.Client
	StopProcessing     int
	GroupProcessing    int
	ContinueProcessing int
}

func (d *jsHome) Plugin(name string) jsPlugin {

	for val, rpc := range d.pluginList {
		if val == name {
			fmt.Println(">> setting plugin", val)
			return jsPlugin{
				client: rpc,
				name:   name,
			}
		}
	}

	return jsPlugin{}
}

func (d *jsHome) GetDeviceByName(name string) jsDevice {
	// fmt.Println("hello")
	for _, v := range d.devices {
		if v.Name == name {
			return v
		}
	}

	return jsDevice{}
}

func (d *jsHome) GetDevices() []jsDevice {
	var out []jsDevice
	// fmt.Println("hello")
	for _, v := range d.devices {
		out = append(out, v)
	}

	return out
}

func (d *jsHome) GetDevicesStartName(s string) []jsDevice {
	var out []jsDevice
	// fmt.Println("hello")
	for k, v := range d.devices {
		if strings.HasPrefix(k, s) {
			out = append(out, v)
		}
	}

	return out
}

func (d *jsHome) GetGroupByName(s string) []jsDevice {
	var out []jsDevice
	// fmt.Println("hello")
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
