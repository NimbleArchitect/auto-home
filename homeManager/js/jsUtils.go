package js

import (
	"fmt"
	"log"
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

	// r.deviceCode[n] = object.(*goja.Object)
}

// processOnTrigger call processes the properties and call the *_ontrigger for each property
//
// dev is then updated with the new properties and values
func (r *JavascriptVM) processOnTrigger(deviceid string, timestamp time.Time, props []map[string]interface{}, dev *jsDevice) {

	for _, prop := range props {
		rawName, ok := prop["name"]
		if !ok {
			log.Println("recieved property without a name")
			continue
		}
		name := rawName.(string)
		if val, ok := prop["type"]; ok {
			log.Printf("processing %s property: %s", val.(string), name)
			switch val.(string) {
			case "switch":
				oldValue := r.deviceState[deviceid].propSwitch[name].state

				swi, err := mapToJsSwitch(prop)
				if err != nil {
					log.Println(err)
				} else {
					_, err := r.RunJS(deviceid, BuildOnAction(name, StrOnTrigger), r.runtime.ToValue(swi.Value))
					if err != nil {
						log.Println(err)
					}

					if oldValue != swi.state {
						dev.propSwitch[name] = swi
					}
					r.deviceState[deviceid].propSwitch[name] = swi
				}

			case "dial":
				oldValue := r.deviceState[deviceid].propDial[name]

				dial, err := mapToJsDial(prop)
				if err != nil {
					log.Println(err)
				} else {
					// check min and max are within range
					if dial.Value > oldValue.max {
						dial.Value = oldValue.max
					}
					if dial.Value < oldValue.min {
						dial.Value = oldValue.min
					}
					_, err := r.RunJS(deviceid, BuildOnAction(name, StrOnTrigger), r.runtime.ToValue(dial.Value))
					if err != nil {
						log.Println(err)
					}

					if oldValue.Value != dial.Value {
						dev.propDial[name] = dial
					}
					r.deviceState[deviceid].propDial[name] = dial
				}

			case "button":
				oldValue := r.deviceState[deviceid].propButton[name].Value

				button, err := mapToJsButton(prop)
				if err != nil {
					log.Println(err)
				} else {
					_, err := r.RunJS(deviceid, BuildOnAction(name, StrOnTrigger), r.runtime.ToValue(button.Value))
					if err != nil {
						log.Println(err)
					}

					if oldValue != button.Value {
						dev.propButton[name] = button
					}
					r.deviceState[deviceid].propButton[name] = button
				}

			case "text":
				oldValue := r.deviceState[deviceid].propText[name].Value

				text, err := mapToJsText(prop)
				if err != nil {
					log.Println(err)
				} else {
					_, err := r.RunJS(deviceid, BuildOnAction(name, StrOnTrigger), r.runtime.ToValue(text.Value))
					if err != nil {
						log.Println(err)
					}

					if oldValue != text.Value {
						dev.propText[name] = text
					}
					r.deviceState[deviceid].propText[name] = text
				}

			default:
				fmt.Println("unknown property type")
			}
		}
	}
}

// processOnChange call loops through all the properties and call the *_onchange for each property
//
// once _onchange has been called the changed value is sent to Updater.Update*
func (r *JavascriptVM) processOnChange(deviceid string, timestamp time.Time, props []map[string]interface{}, dev *jsDevice) {

	for name, swi := range dev.propSwitch {
		// all state props have been updated for the device so we call onchange with the property that was changed
		_, err := r.RunJS(deviceid, BuildOnAction(name, StrOnChange), r.runtime.ToValue(swi.Value))
		if err != nil {
			log.Println(err)
		}
		// now everything has finished we can update the device props
		// save value to device state
		err = r.Updater.UpdateSwitch(deviceid, name, swi.Value)
		if err != nil {
			log.Println("unable to update device state:", err)
		}
	}

	for name, dial := range dev.propDial {
		_, err := r.RunJS(deviceid, BuildOnAction(name, StrOnChange), r.runtime.ToValue(dial.Value))
		if err != nil {
			log.Println(err)
		}
		// save value to device state
		err = r.Updater.UpdateDial(deviceid, name, dial.Value)
		if err != nil {
			log.Println("unable to update device state:", err)
		}
	}

	for name, but := range dev.propButton {
		// all state props have been updated for the device so we call onchange with the property that was changed
		_, err := r.RunJS(deviceid, BuildOnAction(name, StrOnChange), r.runtime.ToValue(but.label))
		if err != nil {
			log.Println(err)
		}
		// now everything has finished we can update the device props
		// save value to device state
		err = r.Updater.UpdateButton(deviceid, name, but.label)
		if err != nil {
			log.Println("unable to update device state:", err)
		}
	}

	for name, txt := range dev.propText {
		// all state props have been updated for the device so we call onchange with the property that was changed
		_, err := r.RunJS(deviceid, BuildOnAction(name, StrOnChange), r.runtime.ToValue(txt.Value))
		if err != nil {
			log.Println(err)
		}
		// now everything has finished we can update the device props
		// save value to device state
		err = r.Updater.UpdateText(deviceid, name, txt.Value)
		if err != nil {
			log.Println("unable to update device state:", err)
		}
	}
}

func (r *JavascriptVM) processGroupChange(deviceid string, timestamp time.Time, props []map[string]interface{}, dev *jsDevice) int {
	var finisheAfterGroups bool

	for _, group := range r.groups {
		for _, v := range group.devices {
			// fmt.Println("6>>", v, deviceid)
			if v == deviceid {
				// TODO: how to I run the group scrip functions??
				val, err := r.RunJSGroupAction(group.Id, "onchange", props)
				if err != nil {
					log.Println(err)
				} else {
					// r.runtime.ToValue(val).ToInteger()
					if val == nil {
						continue
					}

					continueFlag := r.runtime.ToValue(val).ToInteger()
					if continueFlag == FLAG_STOPPROCESSING {
						return int(continueFlag)
					}
					if continueFlag == FLAG_GROUPPROCESSING {
						finisheAfterGroups = true
					}
				}
			}
		}
	}

	if finisheAfterGroups {
		return FLAG_STOPPROCESSING
	}

	return FLAG_CONTINUEPROCESSING
}
