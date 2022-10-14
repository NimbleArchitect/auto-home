package deviceManager

import (
	"log"
	"server/booltype"
	"sync"
	"time"
)

type Button struct {
	lock sync.RWMutex
	data ButtonProperty
}
type ButtonProperty struct {
	Id                    string
	Name                  string
	Description           string
	Value                 booltype.BoolType
	Previous              bool
	Mode                  uint
	repeatWindowTimeStamp time.Time
	repeatWindowDuration  time.Duration
}

func (d *Device) NewButton() *Button {
	return &Button{
		lock: sync.RWMutex{},
		data: ButtonProperty{},
	}
}

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

func (d *Device) SetButton(name string, property *ButtonProperty) {
	prop, ok := d.PropertyButton[name]
	if !ok {
		d.PropertyButton[name] = &Button{
			lock: sync.RWMutex{},
			data: *property,
		}
	} else {
		prop.lock.Lock()
		prop.data = *property
		prop.lock.Unlock()
	}
}

func (d *Device) ButtonValue(name string) (string, bool) {
	property, ok := d.PropertyButton[name]
	if ok {
		property.lock.RLock()
		data := property.data
		property.lock.RUnlock()
		return data.Value.String(), true
	}

	return "", false
}

// Was UpdateButton
func (d *Device) SetButtonValue(name string, value string) {
	property, ok := d.PropertyButton[name]
	if ok {
		property.lock.Lock()
		property.data.Value.Set(value)
		property.lock.Unlock()
	}

}

func (d *Device) ButtonWindow(name string, timestamp time.Time) bool {
	property, ok := d.PropertyButton[name]
	if ok {
		data := property.data

		if data.repeatWindowTimeStamp.Before(timestamp) {
			// fmt.Println("d>>", data.repeatWindowDuration.Milliseconds())
			newExpire := timestamp.Add(data.repeatWindowDuration)

			property.lock.Lock()
			property.data.repeatWindowTimeStamp = newExpire
			property.lock.Unlock()

			// fmt.Println("**>> update allowed")
			return true
		} else {
			// fmt.Println("**>> update blocked")
			return false
		}
	}

	return false
}

func (d *Device) SetButtonWindow(name string, duration int64) {
	property, ok := d.PropertyButton[name]
	if ok {
		property.lock.Lock()
		property.data.repeatWindowDuration = time.Duration(duration) * time.Millisecond
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
		return nil, ErrMissingPropertyValue
	} else {
		prop.Value.Set(v.(string))
	}

	if v, ok := props["mode"]; !ok {
		return nil, ErrMissingPropertyMode
	} else {
		prop.Mode, err = GetModeFromString(v.(string))
		if err != nil {
			log.Println(err)
		}
	}

	return &prop, nil

}
