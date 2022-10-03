package js

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dop251/goja"
)

type DeviceUpdator interface {
	UpdateDial(string, string, int) error
	UpdateSwitch(string, string, string) error
	UpdateButton(string, string, string) error
	UpdateText(string, string, string) error

	GetDialValue(string, string) (int, bool)
	GetSwitchValue(string, string) (string, bool)
	GetButtonValue(string, string) (string, bool)
	GetTextValue(string, string) (string, bool)

	// GetDialHistory()

	RunGroupAction(string, string, []map[string]interface{}) (interface{}, error)
}

type JavascriptVM struct {
	runtime     *goja.Runtime
	deviceState map[string]jsDevice
	groups      map[string]jsGroup
	Updater     DeviceUpdator
}

func NewScript(actionFile string) (*JavascriptVM, error) {
	var js JavascriptVM

	vm := goja.New()
	// vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	vm.SetFieldNameMapper(goja.UncapFieldNameMapper())

	var console jsConsole
	err := vm.Set("console", console)
	if err != nil {
		log.Println(err)
	}

	vm.Set("thread", runAsThread)
	if err != nil {
		log.Println(err)
	}

	// vm.SetFieldNameMapper(goja.UncapFieldNameMapper())

	file, err := os.ReadFile(actionFile)
	if err != nil {
		log.Println("unable to read file:", err)
	}

	actionScript, err := goja.Compile(actionFile, string(file), true)
	if err != nil {
		log.Println("unable to compile script", err)
	}

	_, err = vm.RunProgram(actionScript)
	if err != nil {
		log.Println("unable to run script", actionScript, "error reported was", err)
	}

	js.runtime = vm
	return &js, nil

}

// runAsThread runs the js function as a new thread, this could be dangerous/not thread safe
func runAsThread(obj goja.Value, val goja.Value) {
	call, ok := goja.AssertFunction(obj)
	if ok {
		go call(goja.Undefined(), val)
	} else {
		log.Println("thread call not a function")
	}
}

func (r *JavascriptVM) RunJSGroupAction(fnName string, props []map[string]interface{}) (interface{}, error) {
	// var dev jsDevice

	log.Println("event triggered")

	// dev.propSwitch = make(map[string]jsSwitch)
	// dev.propDial = make(map[string]jsDial)
	// dev.propButton = make(map[string]jsButton)
	// dev.propText = make(map[string]jsText)

	// log.Println("state:", m.devices)

	// lookup changes and trigger change notifications
	out, err := r.RunJS(fnName, r.runtime.ToValue(props))
	if err != nil {
		log.Println(err)
	}

	return out, nil
}

func (r *JavascriptVM) Process(deviceid string, timestamp time.Time, props []map[string]interface{}) {
	var dev jsDevice

	log.Println("event triggered")

	dev.propSwitch = make(map[string]jsSwitch)
	dev.propDial = make(map[string]jsDial)
	dev.propButton = make(map[string]jsButton)
	dev.propText = make(map[string]jsText)

	// log.Println("state:", m.devices)

	// lookup changes and trigger change notifications
	r.processOnTrigger(deviceid, timestamp, props, &dev)

	// TODO: not sure this is the correct order as it depends on if we wnat groups to return a no further processing argument
	continueFlag := r.processGroupChange(deviceid, timestamp, props, &dev)
	if continueFlag != FLAG_STOPPROCESSING {
		r.processOnChange(deviceid, timestamp, props, &dev)
	}

}
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
				oldValue := r.deviceState[deviceid].propSwitch[name].Value

				swi, err := mapToJsSwitch(prop)
				if err != nil {
					log.Println(err)
				} else {
					_, err := r.RunJS(name+"_ontrigger", r.runtime.ToValue(swi.label))
					if err != nil {
						log.Println(err)
					}

					if oldValue != swi.Value {
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
					_, err := r.RunJS(name+"_ontrigger", r.runtime.ToValue(dial.Value))
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
					_, err := r.RunJS(name+"_ontrigger", r.runtime.ToValue(button.Value))
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
					_, err := r.RunJS(name+"_ontrigger", r.runtime.ToValue(text.Value))
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

func (r *JavascriptVM) processOnChange(deviceid string, timestamp time.Time, props []map[string]interface{}, dev *jsDevice) {

	for name, swi := range dev.propSwitch {
		// all state props have been updated for the device so we call onchange with the property that was changed
		_, err := r.RunJS(name+"_onchange", r.runtime.ToValue(swi.label))
		if err != nil {
			log.Println(err)
		}
		// now everything has finished we can update the device props
		// save value to device state
		err = r.Updater.UpdateSwitch(deviceid, name, swi.label)
		if err != nil {
			log.Println("unable to update device state:", err)
		}
	}

	for name, dial := range dev.propDial {
		_, err := r.RunJS(name+"_onchange", r.runtime.ToValue(dial.Value))
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
		_, err := r.RunJS(name+"_onchange", r.runtime.ToValue(but.label))
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
		_, err := r.RunJS(name+"_onchange", r.runtime.ToValue(txt.Value))
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
			if v == deviceid {
				// TODO: how to I run the group scrip functions??
				val, err := r.Updater.RunGroupAction(group.Id, "group_onchange", props)
				if err != nil {
					log.Println(err)
				} else {
					// r.runtime.ToValue(val).ToInteger()
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
