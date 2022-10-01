package js

import (
	"strings"
	"time"
)

type jsHome struct {
	devices map[string]jsDevice
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

func (d *jsHome) Sleep(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}
