package js

import (
	"log"
	"server/deviceManager"
	"strings"
	"time"

	"github.com/dop251/goja"
)

func (r *JavascriptVM) objLoader(name goja.Value, object goja.Value) {
	// TODO: make name safe as its use input
	n := name.String()

	parts := strings.Split(n, "/")
	switch parts[0] {
	case "group":
		r.groupCode[n] = object.(*goja.Object)
	case "user":
		r.userCode[n] = object.(*goja.Object)
	case "device":
		fallthrough
	default:
		r.deviceCode[n] = object.(*goja.Object)
	}
}

// processOnTrigger call processes the properties and call the *_ontrigger for each property
//
// dev is then updated with the new properties and values
func (r *JavascriptVM) processOnTrigger(deviceid string, timestamp time.Time, props JSPropsList, dev *jsDevice) {
	_, ok := r.deviceState[deviceid]
	if !ok {
		log.Println("processTrigger device not found", deviceid)
		return
	}

	for name, swi := range props.propSwitch {
		oldValue := r.deviceState[deviceid].propSwitch[name].state
		val, err := r.RunJS(deviceid, BuildOnAction(name, StrOnTrigger), r.runtime.ToValue(swi.Value))
		if err != nil {
			log.Println(err)
		} else {
			if val != nil {
				swi.flag.Set(r.runtime.ToValue(val).ToInteger())
			}
		}

		if oldValue != swi.state {
			dev.propSwitch[name] = swi
		}
		r.deviceState[deviceid].propSwitch[name] = swi
	}

	for name, dial := range props.propDial {
		oldValue := r.deviceState[deviceid].propDial[name]

		// check min and max are within range
		if dial.Value > oldValue.max {
			dial.Value = oldValue.max
		}
		if dial.Value < oldValue.min {
			dial.Value = oldValue.min
		}
		val, err := r.RunJS(deviceid, BuildOnAction(name, StrOnTrigger), r.runtime.ToValue(dial.Value))
		if err != nil {
			log.Println(err)
		} else {
			if val != nil {
				dial.flag.Set(r.runtime.ToValue(val).ToInteger())
			}
		}

		if oldValue.Value != dial.Value {
			dev.propDial[name] = dial
		}
		r.deviceState[deviceid].propDial[name] = dial
	}

	for name, button := range props.propButton {
		oldValue := r.deviceState[deviceid].propButton[name].Value

		val, err := r.RunJS(deviceid, BuildOnAction(name, StrOnTrigger), r.runtime.ToValue(button.Value))
		if err != nil {
			log.Println(err)
		} else {
			if val != nil {
				button.flag.Set(r.runtime.ToValue(val).ToInteger())
			}
		}

		if oldValue != button.Value {
			dev.propButton[name] = button
		}
		r.deviceState[deviceid].propButton[name] = button
	}

	for name, text := range props.propText {
		oldValue := r.deviceState[deviceid].propText[name].Value

		_, err := r.RunJS(deviceid, BuildOnAction(name, StrOnTrigger), r.runtime.ToValue(text.Value))
		if err != nil {
			log.Println(err)
		}

		if oldValue != text.Value {
			dev.propText[name] = text
		}
		r.deviceState[deviceid].propText[name] = text
	}
}

// processOnChange call loops through all the properties and call the *_onchange for each property
//
// once _onchange has been called the changed value is sent to Updater.Update*
func (r *JavascriptVM) processOnChange(deviceid string, dev *jsDevice, FLAG int) {
	var liveDevice *deviceManager.Device
	tmp, ok := r.deviceState[deviceid]
	if ok {
		liveDevice = tmp.liveDevice
	}

	for name, swi := range dev.propSwitch {
		if FLAG != FLAG_STOPPROCESSING {
			// all state props have been updated for the device so we call onchange with the property that was changed
			val, err := r.RunJS(deviceid, BuildOnAction(name, StrOnChange), r.runtime.ToValue(swi.Value))
			if err != nil {
				log.Println(err)
			} else {
				if val != nil {
					swi.flag.Set(r.runtime.ToValue(val).ToInteger())
				}
			}
		}

		// now everything has finished we can update the device props
		// save value to device state

		if liveDevice != nil && swi.flag.Not(FLAG_PREVENTUPDATE) {
			liveDevice.SetSwitchValue(name, swi.Value)
		}
		// if err != nil {
		// 	log.Println("unable to update device state:", err)
		// }
	}

	for name, dial := range dev.propDial {
		if FLAG != FLAG_STOPPROCESSING {
			val, err := r.RunJS(deviceid, BuildOnAction(name, StrOnChange), r.runtime.ToValue(dial.Value))
			if err != nil {
				log.Println(err)
			} else {
				if val != nil {
					dial.flag.Set(r.runtime.ToValue(val).ToInteger())
				}
			}
		}
		// save value to live device
		if liveDevice != nil && dial.flag.Not(FLAG_PREVENTUPDATE) {
			liveDevice.SetDialValue(name, dial.Value)
		}
		// err = r.Updater.UpdateDial(deviceid, name, dial.Value)
		// if err != nil {
		// 	log.Println("unable to update device state:", err)
		// }
	}

	for name, button := range dev.propButton {
		if FLAG != FLAG_STOPPROCESSING {
			// all state props have been updated for the device so we call onchange with the property that was changed
			val, err := r.RunJS(deviceid, BuildOnAction(name, StrOnChange), r.runtime.ToValue(button.label))
			if err != nil {
				log.Println(err)
			} else {
				if val != nil {
					button.flag.Set(r.runtime.ToValue(val).ToInteger())
				}
			}
		}
		// now everything has finished we can update the device props
		// save value to device state
		if liveDevice != nil && button.flag.Not(FLAG_PREVENTUPDATE) {
			liveDevice.SetButtonValue(name, button.label)
		}
		// if err != nil {
		// 	log.Println("unable to update device state:", err)
		// }
	}

	for name, txt := range dev.propText {
		if FLAG != FLAG_STOPPROCESSING {
			// all state props have been updated for the device so we call onchange with the property that was changed
			val, err := r.RunJS(deviceid, BuildOnAction(name, StrOnChange), r.runtime.ToValue(txt.Value))
			if err != nil {
				log.Println(err)
			} else {
				if val != nil {
					txt.flag.Set(r.runtime.ToValue(val).ToInteger())
				}
			}
		}
		// now everything has finished we can update the device props
		// save value to device state
		if liveDevice != nil && txt.flag.Not(FLAG_PREVENTUPDATE) {
			liveDevice.SetTextValue(name, txt.Value)
		}
		// if err != nil {
		// 	log.Println("unable to update device state:", err)
		// }
	}
}

func (r *JavascriptVM) processGroupChange(deviceid string, props JSPropsList) int {
	var finishAfterGroups bool
	var searchList []jsGroup

	for _, group := range r.groups {
		for _, v := range group.devices {
			if v == deviceid {
				searchList = append(searchList, group)
			}
		}
	}

	jsRunList := make(map[string]jsGroup)

	for i := 0; i < 60; i++ {
		var newList []jsGroup
		var carryOn bool

		// build a list of unique parent groups so we can run any attached scripts later
		for _, group := range searchList {
			parents := r.ParentsOf(group.Id)
			if len(parents) <= 0 {
				parents[group.Id] = group
			}
			for _, v := range parents {
				if _, ok := jsRunList[v.Id]; !ok {
					newList = append(newList, v)
					carryOn = true
				}
				jsRunList[v.Id] = v
			}

		}

		if carryOn {
			searchList = newList
		} else {
			break
		}
	}

	// TODO: need to call jsRun for every group from closest to device to furthest
	for _, group := range jsRunList {
		if group.liveGroup != nil {
			if group.liveGroup.Window(time.Now()) {
				continue
			}
		}

		continueFlag, err := r.RunJSGroup(group.Id, props)
		if err != nil {
			log.Println("group error", err)
			continue
		}

		if continueFlag == FLAG_STOPPROCESSING {
			return int(continueFlag)
		}
		if continueFlag == FLAG_GROUPPROCESSING {
			finishAfterGroups = true
		}
	}

	if finishAfterGroups {
		return FLAG_STOPPROCESSING
	}

	return FLAG_CONTINUEPROCESSING
}

func (r *JavascriptVM) ParentsOf(name string) map[string]jsGroup {
	foundMap := make(map[string]jsGroup)

	// TODO: this needs to be recursive
	for _, group := range r.groups {
		for _, child := range group.groups {
			if child == name {
				foundMap[group.Id] = group
			}
		}
	}

	return foundMap
}
