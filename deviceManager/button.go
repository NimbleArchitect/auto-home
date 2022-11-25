package deviceManager

import (
	"fmt"
	"log"
	"server/booltype"
	"sync"
	"time"
)

type Button struct {
	lock    sync.RWMutex
	data    ButtonProperty
	history struct {
		max    int
		index  int
		values []bool
	}
	repeatWindowTimeStamp time.Time
	repeatWindowDuration  time.Duration
}
type ButtonProperty struct {
	Id          string
	Name        string
	Description string
	Value       booltype.BoolType
	// Previous              bool
	Mode uint
}

// func (d *Device) NewButton() *Button {
// 	return &Button{
// 		lock: sync.RWMutex{},
// 		data: ButtonProperty{},
// 	}
// }

func (d *Device) ButtonKeys() []string {

	return d.ButtonNames
}

func (d *Device) Button(name string) ButtonProperty {
	property, ok := d.PropertyButton[name]
	if ok {
		property.lock.RLock()
		out := property.data
		property.lock.RUnlock()
		return out
	}

	return ButtonProperty{}
}

func (d *Device) ButtonAsMap() map[string]ButtonProperty {
	propertyList := make(map[string]ButtonProperty, len(d.PropertyButton))

	for name, property := range d.PropertyButton {
		property.lock.RLock()
		propertyList[name] = property.data
		property.lock.RUnlock()
	}

	return propertyList
}

func (d *Device) ButtonHistory(name string, value int) (ButtonProperty, bool) {
	// TODO: need finishing should this return a single value or a property

	return ButtonProperty{}, false
}

func (d *Device) SetButton(name string, property *ButtonProperty) {
	prop, ok := d.PropertyButton[name]
	if !ok {
		duration := d.repeatWindow[name]
		log.Println("new button", name, duration)

		button := Button{
			lock:                 sync.RWMutex{},
			data:                 *property,
			repeatWindowDuration: time.Duration(duration) * time.Millisecond,
		}
		button.history.max = d.maxPropertyHistory

		d.PropertyButton[name] = &button
		d.ButtonNames = append(d.ButtonNames, name)
	} else {
		log.Println("overwriting button", name)
		prop.lock.Lock()
		prop.data = *property
		prop.lock.Unlock()
	}
}

func (d *Device) ButtonValue(name string) (booltype.BoolType, bool) {
	property, ok := d.PropertyButton[name]
	if ok {
		property.lock.RLock()
		data := property.data
		property.lock.RUnlock()
		return data.Value, true
	}

	return booltype.BoolType{}, false
}

// updates the live device
func (d *Device) WriteButtonValue(name string, value string) {
	fmt.Println("F:WriteButtonValue:d.Id =", d.Id)

	if d.clientConnection != nil {
		if writer := d.clientConnection.ClientWriter(d.ClientId); writer != nil {
			jsonOut := d.MakeAction(d.Id, name, BUTTON, fmt.Sprint(value))
			writer.Write(jsonOut)
		}
	}
}

// Was UpdateButton
func (d *Device) SetButtonValue(name string, value string) {
	property, ok := d.PropertyButton[name]
	if ok {
		property.lock.Lock()
		if property.history.index >= property.history.max {
			property.history.index = 0
		}

		if property.history.max-1 >= len(property.history.values) {
			property.history.values = append(property.history.values, property.data.Value.GetBool())
		} else {
			property.history.values[property.history.index] = property.data.Value.GetBool()
		}
		property.history.index++

		property.data.Value.Set(value)
		// copy the Id so we can unlock before we start the call back action, this means we dont have to
		//  keep the lock open until the client has rwsponded

		property.lock.Unlock()
	}
}

func (d *Device) ButtonWindow(name string, timestamp time.Time) bool {
	property, ok := d.PropertyButton[name]
	if ok {
		if property.repeatWindowTimeStamp.Before(timestamp) {
			newExpire := timestamp.Add(property.repeatWindowDuration)
			property.lock.Lock()
			property.repeatWindowTimeStamp = newExpire
			property.lock.Unlock()
			return true
		} else {
			return false
		}
	}

	return false
}

func (d *Device) SetButtonWindow(name string, duration int64) {
	property, ok := d.PropertyButton[name]
	if ok {
		timelimit := time.Duration(duration) * time.Millisecond
		property.lock.Lock()
		property.repeatWindowDuration = timelimit
		property.lock.Unlock()
	}
}

// Was ReadPropertyButton
func (d *Device) Map2Button(props map[string]interface{}) (*ButtonProperty, error) {
	var prop ButtonProperty
	var err error

	log.Println("reading button property")
	if v, ok := props["name"]; !ok {
		return nil, ErrMissingPropertyName
	} else {
		// TODO: clean the string
		prop.Name = v.(string)
		log.Println("name", prop.Name)
	}

	if v, ok := props["description"]; ok {
		// TODO: clean the string
		prop.Description = v.(string)
	}

	if v, ok := props["value"]; !ok {
		fmt.Println("error missing value")
		return nil, ErrMissingPropertyValue
	} else {
		prop.Value.SetBool(v.(bool))
		// prop.Value.Set(v.(string))
	}

	if v, ok := props["mode"]; !ok {
		fmt.Println("error mssing mode")
		return nil, ErrMissingPropertyMode
	} else {
		prop.Mode, err = GetModeFromString(v.(string))

		if err != nil {
			log.Println(err)
		}
	}

	return &prop, nil

}
