package js

import (
	"server/deviceManager"
	log "server/logger"
	"strings"
	"time"

	"github.com/dop251/goja"
)

const setNameIgnoreString = "!`¬\"£^*|'~<>" // list of chars that we dont accept in a set name
const setNameSplitCharacter = "/"

type groupObject struct {
	Id       string
	Name     string
	Property interface{}
}

// objLoader entry point for javascript set function
func (r *JavascriptVM) objLoader(name goja.Value, object goja.Value) {
	n := name.String()
	if strings.ContainsAny(n, setNameIgnoreString) {
		log.Errorf("Invalid charactors detected in set name, the following charactors are invalid: \"%s\"", setNameIgnoreString)
		return
	}

	parts := strings.Split(n, setNameSplitCharacter)
	switch parts[0] {
	case "plugin":
		r.pluginCode[n] = object.(*goja.Object)
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
// TODO: needs merging with processGroupChange as the ontrigger event is being removed
func (r *JavascriptVM) processOnTrigger(deviceid string, timestamp time.Time, props JSPropsList, dev *jsDevice) []interface{} {
	var propsArray []interface{}

	_, ok := r.deviceState[deviceid]
	if !ok {
		log.Error("processTrigger device not found", deviceid)
		return nil
	}

	for name, swi := range props.propSwitch {
		oldValue := r.deviceState[deviceid].propSwitch[name].state
		// val, err := r.RunJS(deviceid, BuildOnAction(name, StrOnTrigger), r.runtime.ToValue(swi.Value))
		// if err != nil {
		// 	log.Error(err)
		// } else {
		// 	if val != nil {
		// 		swi.flag.Set(r.runtime.ToValue(val).ToInteger())
		// 	}
		// }

		if oldValue != swi.state {
			dev.propSwitch[name] = swi
			propsArray = append(propsArray, swi)
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
		// val, err := r.RunJS(deviceid, BuildOnAction(name, StrOnTrigger), r.runtime.ToValue(dial.Value))
		// if err != nil {
		// 	log.Error(err)
		// } else {
		// 	if val != nil {
		// 		dial.flag.Set(r.runtime.ToValue(val).ToInteger())
		// 	}
		// }

		if oldValue.Value != dial.Value {
			dev.propDial[name] = dial
			propsArray = append(propsArray, dial)
		}
		r.deviceState[deviceid].propDial[name] = dial
	}

	for name, button := range props.propButton {
		oldValue := r.deviceState[deviceid].propButton[name].Value

		// val, err := r.RunJS(deviceid, BuildOnAction(name, StrOnTrigger), r.runtime.ToValue(button.Value))
		// if err != nil {
		// 	log.Error(err)
		// } else {
		// 	if val != nil {
		// 		button.flag.Set(r.runtime.ToValue(val).ToInteger())
		// 	}
		// }

		if oldValue != button.Value {
			dev.propButton[name] = button
			propsArray = append(propsArray, button)
		}
		r.deviceState[deviceid].propButton[name] = button
	}

	for name, text := range props.propText {
		oldValue := r.deviceState[deviceid].propText[name].Value

		// _, err := r.RunJS(deviceid, BuildOnAction(name, StrOnTrigger), r.runtime.ToValue(text.Value))
		// if err != nil {
		// 	log.Error(err)
		// }

		if oldValue != text.Value {
			dev.propText[name] = text
			propsArray = append(propsArray, text)
		}
		r.deviceState[deviceid].propText[name] = text
	}

	return propsArray
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
		// now everything has finished we can update the device props
		// save value to device state

		if liveDevice != nil && swi.flag.Not(FLAG_PREVENTUPDATE) {
			liveDevice.SetSwitchValue(name, swi.Value)
		}

		if FLAG != FLAG_STOPPROCESSING {
			// all state props have been updated for the device so we call onchange with the property that was changed
			val, err := r.RunJS(deviceid, BuildOnAction(name, StrOnChange), r.runtime.ToValue(swi.Value))
			if err != nil {
				log.Error(err)
			} else {
				if val != nil {
					swi.flag.Set(r.runtime.ToValue(val).ToInteger())
				}
			}
		}
	}

	for name, dial := range dev.propDial {
		// save value to live device
		if liveDevice != nil && dial.flag.Not(FLAG_PREVENTUPDATE) {
			liveDevice.SetDialValue(name, dial.Value)
		}

		if FLAG != FLAG_STOPPROCESSING {
			val, err := r.RunJS(deviceid, BuildOnAction(name, StrOnChange), r.runtime.ToValue(dial.Value))
			if err != nil {
				log.Error(err)
			} else {
				if val != nil {
					dial.flag.Set(r.runtime.ToValue(val).ToInteger())
				}
			}
		}
	}

	for name, button := range dev.propButton {
		// now everything has finished we can update the device props
		// save value to device state
		if liveDevice != nil && button.flag.Not(FLAG_PREVENTUPDATE) {
			liveDevice.SetButtonValue(name, button.Value)
		}

		if FLAG != FLAG_STOPPROCESSING {
			// all state props have been updated for the device so we call onchange with the property that was changed
			val, err := r.RunJS(deviceid, BuildOnAction(name, StrOnChange), r.runtime.ToValue(button.Value))
			if err != nil {
				log.Error(err)
			} else {
				if val != nil {
					button.flag.Set(r.runtime.ToValue(val).ToInteger())
				}
			}
		}
	}

	for name, txt := range dev.propText {
		// now everything has finished we can update the device props
		// save value to device state
		if liveDevice != nil && txt.flag.Not(FLAG_PREVENTUPDATE) {
			liveDevice.SetTextValue(name, txt.Value)
		}

		if FLAG != FLAG_STOPPROCESSING {
			// all state props have been updated for the device so we call onchange with the property that was changed
			val, err := r.RunJS(deviceid, BuildOnAction(name, StrOnChange), r.runtime.ToValue(txt.Value))
			if err != nil {
				log.Error(err)
			} else {
				if val != nil {
					txt.flag.Set(r.runtime.ToValue(val).ToInteger())
				}
			}
		}
	}
}

func (r *JavascriptVM) processGroupChange(deviceid string, props []interface{}) int {
	var finishAfterGroups bool
	var searchList []jsGroup
	var deviceName string

	// first get a list of groups that have our device as a member
	for _, group := range r.groups {
		for _, v := range group.devices {
			if v == deviceid {
				searchList = append(searchList, group)
			}
		}
	}

	jsRunList := make(map[string]jsGroup)

	// 60 is an upper limit on the number of groups we will attempt to check, we should bail before
	//  hitting this limit so this is just a safety limit
	for i := 0; i < 60; i++ {
		var newList []jsGroup
		var carryOn bool

		// build a list of unique parent groups so we can run any attached scripts later
		// this removes any duplicates so we only run a group once
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

	if value, ok := r.deviceState[deviceid]; ok {
		deviceName = value.Name
	}
	// great a new group properties object to pass over to RunJsGroup()
	// this sends the deviceID, name and supplied properties that have changed
	groupObj := groupObject{
		Id:       deviceid,
		Name:     deviceName,
		Property: props,
	}

	// TODO: need to call jsRun for every group from closest to device to furthest
	for _, group := range jsRunList {
		if group.liveGroup != nil {
			if group.liveGroup.Window(time.Now()) {
				continue
			}
		}

		// expose groupObj to javascript
		continueFlag, err := r.RunJSGroup(group.Id, groupObj)
		if err != nil {
			log.Error("group error", err)
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

// ParentsOf list all parent groups of the group called name
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

// setJsGlobal sets up the objects for the runtime, these are the object that can
//
//	change so need loading everytime the vm is called, currently sets up home and plugin obects
func (r *JavascriptVM) setJsGlobal() jsHome {

	home := jsHome{
		vm:      r,
		devices: r.deviceState,

		StopProcessing:     FLAG_STOPPROCESSING,
		ContinueProcessing: FLAG_CONTINUEPROCESSING,
		GroupProcessing:    FLAG_GROUPPROCESSING,
		PreventUpdate:      FLAG_PREVENTUPDATE,
	}

	if err := r.runtime.Set("plugin", r.plugins); err != nil {
		log.Error("unable to attach plugin object to javascript vm:", err)
	}
	if err := r.runtime.Set("home", home); err != nil {
		log.Error("unable to attach plugin object to javascript vm:", err)
	}

	return home
}
