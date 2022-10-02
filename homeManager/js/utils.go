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
	UpdateButton(string, string, bool) error
	UpdateText(string, string, string) error
	GetDialValue(string, string) (int, bool)
	GetSwitchValue(string, string) (string, bool)
	GetButtonValue(string, string) (bool, bool)
	GetTextValue(string, string) (string, bool)
	//TODO: add button and text props

}

type JavascriptVM struct {
	runtime     *goja.Runtime
	deviceState map[string]jsDevice
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

	vm.SetFieldNameMapper(goja.UncapFieldNameMapper())

	// TODO: needs fixing, can seem to use goja to read the file I have to do it by hand :(
	file, err := os.ReadFile(actionFile)

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

func (r *JavascriptVM) Process(deviceid string, timestamp time.Time, props []map[string]interface{}) {
	var dev jsDevice

	log.Println("event triggered")

	if len(dev.propSwitch) == 0 {
		dev.propSwitch = make(map[string]jsSwitch)
	}

	if len(dev.propDial) == 0 {
		dev.propDial = make(map[string]jsDial)
	}

	if len(dev.propButton) == 0 {
		dev.propButton = make(map[string]jsButton)
	}

	if len(dev.propText) == 0 {
		dev.propText = make(map[string]jsText)
	}

	// log.Println("state:", m.devices)
	// save the current state of all devices
	// jsState, _ := m.SaveState()

	// lookup changes, trigger change notifications, what am I supposed
	//  to trigger and how am I supposed to trigger it???

	// lookup device, trigger device scripts
	// dev := m.devices[deviceid]
	// fmt.Println(">>", deviceid)
	// vm := m.actions[deviceid].jsvm

	changeList := map[string]int{}

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
						changeList[name] = 0
					}
					// fmt.Println("3>>", deviceid)
					// fmt.Println("4>>", r.deviceState)

					r.deviceState[deviceid].propSwitch[name] = swi
				}

			case "dial":
				// TODO: check min and max are within range
				oldValue := r.deviceState[deviceid].propDial[name].Value

				dial, err := mapToJsDial(prop)
				if err != nil {
					log.Println(err)
				} else {
					_, err := r.RunJS(name+"_ontrigger", r.runtime.ToValue(dial.Value))
					if err != nil {
						log.Println(err)
					}

					if oldValue != dial.Value {
						dev.propDial[name] = dial
						changeList[name] = 0
					}
					r.deviceState[deviceid].propDial[name] = dial
				}

			case "button":
				// TODO: check min and max are within range
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
						changeList[name] = 0
					}
					// r.deviceState[deviceid].propButton[name] = button
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
						changeList[name] = 0
					}
					r.deviceState[deviceid].propText[name] = text
				}

			default:
				fmt.Println("unknown property type")
			}
		}
	}

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
		_, err := r.RunJS(name+"_onchange", r.runtime.ToValue(but.Value))
		if err != nil {
			log.Println(err)
		}
		// buttons are never updated so we end here
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
