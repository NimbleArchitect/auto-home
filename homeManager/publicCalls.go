package home

import (
	"encoding/json"
	"fmt"
	log "server/logger"
	"strconv"

	"github.com/dop251/goja"
)

// WebCallPlugin calls the function callNAme of the plugin named pluginName using the postData as the arguments,
// only to be called from web interfaces
func (m *Manager) WebCallPlugin(pluginName string, callName string, postData map[string]interface{}) []byte {
	var out map[string]interface{}

	if len(pluginName) <= 0 || len(callName) <= 0 {
		return []byte{}
	}

	if plugin := m.plugins.Get(pluginName); plugin != nil {
		if caller := plugin.Get(callName); caller != nil {
			if len(postData) > 0 {
				v := goja.New().ToValue(postData)
				out = caller.Run(v.Export())
			} else {
				var empty goja.Value
				out = caller.Run(empty)
			}

			mapLen := len(out)
			fmt.Println("mapLen:", mapLen)
			if mapLen == 0 {
				fmt.Println("return empty")
				return []byte{}
			} else if mapLen < 2 {
				fmt.Println("return 1")
				data, _ := json.Marshal(out["0"])
				fmt.Println("return data:", string(data))
				return data
			} else {
				fmt.Println("return all")
				data, _ := json.Marshal(out)
				fmt.Println("return data:", string(data))
				return data
			}
			// data, _ := json.Marshal(out)
			// return data
		}
	}
	return []byte{}
}

// WebAllDeviceAsJson returns a byte array containing a json list of all currently known devices and their properties,
// only to be called from web interfaces
func (m *Manager) WebAllDeviceAsJson() []byte {
	var jsonDeviceList []byte

	jsonDeviceList = append(jsonDeviceList, []byte("[")...)
	deviceList := m.devices.Iterate()
	for deviceList.Next() {
		_, device := deviceList.Get()
		bytesout, err := device.AsJson()
		if err != nil {
			log.Error("unable to convert device to json:", err)
		} else {
			jsonDeviceList = append(jsonDeviceList, append(bytesout, []byte(",")...)...)
		}
	}

	if len(jsonDeviceList) > 1 {
		return append(jsonDeviceList[:len(jsonDeviceList)-1], []byte("]")...)
	}

	return []byte("[]")
}

// WebDeviceAsJson searches for a device matching the provided id and returns the json string in bytes representing the device
func (m *Manager) WebDeviceAsJson(id string) []byte {

	device, ok := m.devices.Device(id)
	if ok {
		jsonDevice, err := device.AsJson()
		if err != nil {
			log.Error("unable to convert device to json:", err)
		} else {
			return jsonDevice
		}
	}

	return []byte{}
}

// WebDevicePropertyAsJson searches for a device and property matching the provided deviceid and property name and returns the property value
func (m *Manager) WebDevicePropertyAsJson(deviceid string, propertyName string) []byte {

	device, ok := m.devices.Device(deviceid)
	if ok {
		if val, found := device.ButtonValue(propertyName); found {
			bytesout := fmt.Sprintf(`{"type":"button", "state":%s, "value":"%t"}`, val.String(), val.Bool())
			return []byte(bytesout)
		}

		if val, found := device.DialValue(propertyName); found {
			bytesout := fmt.Sprintf(`{"type":"dial", "state":%d}`, val)
			return []byte(bytesout)
		}

		if val, found := device.SwitchValue(propertyName); found {
			bytesout := fmt.Sprintf(`{"type":"switch", "state":%s, "value":"%t"}`, val.String(), val.Bool())
			return []byte(bytesout)
		}

		if val, found := device.TextValue(propertyName); found {
			bytesout := fmt.Sprintf(`{"type":"text", "state":%s}`, val)
			return []byte(bytesout)
		}
	}

	return []byte{}
}

// WebSetDeviceProperty set the property propertyName of the device with deviceid to strValue
func (m *Manager) WebSetDeviceProperty(user string, deviceid string, propertyName string, strValue string) bool {
	// TODO: check username has access to make changes to the properties/device
	device, ok := m.devices.Device(deviceid)
	if ok {
		if _, found := device.ButtonValue(propertyName); found {
			device.SetButtonValue(propertyName, strValue)
			return true
		}

		if _, found := device.DialValue(propertyName); found {
			val, err := strconv.Atoi(strValue)
			if err == nil {
				device.SetDialValue(propertyName, val)
				return true
			}
		}

		if _, found := device.SwitchValue(propertyName); found {
			device.SetSwitchValue(propertyName, strValue)
			return true
		}

		if _, found := device.TextValue(propertyName); found {
			device.SetTextValue(propertyName, strValue)
			return true
		}
	}

	return false
}

func (m *Manager) WebSetDevicePropertyPercent(deviceid string, propertyName string, strValue string) bool {
	// TODO: needs finishing
	log.Error("not implemented")
	return false
}
