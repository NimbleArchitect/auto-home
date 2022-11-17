package deviceManager

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Text struct {
	lock    sync.RWMutex
	data    TextProperty
	history struct {
		max    int
		index  int
		values []string
	}
	repeatWindowTimeStamp time.Time
	repeatWindowDuration  time.Duration
}
type TextProperty struct {
	Id          string
	Name        string
	Description string
	Value       string
	// Previous              string
	Mode uint
}

// func (d *Device) NewText() *Text {
// 	return &Text{
// 		lock: sync.RWMutex{},
// 		data: TextProperty{},
// 	}
// }

func (d *Device) TextKeys() []string {

	return d.TextNames
}

func (d *Device) Text(name string) TextProperty {
	property, ok := d.PropertyText[name]
	if ok {
		property.lock.RLock()
		out := property.data
		property.lock.RUnlock()
		return out
	}

	return TextProperty{}
}

func (d *Device) TextAsMap() map[string]TextProperty {
	propertyList := make(map[string]TextProperty, len(d.PropertyText))

	for name, property := range d.PropertyText {
		property.lock.RLock()
		propertyList[name] = property.data
		property.lock.RUnlock()
	}

	return propertyList
}

func (d *Device) SetText(name string, property *TextProperty) {
	prop, ok := d.PropertyText[name]
	if !ok {
		duration := d.repeatWindow[name]
		log.Println("new text", name, duration)

		text := Text{
			lock:                 sync.RWMutex{},
			data:                 *property,
			repeatWindowDuration: time.Duration(duration) * time.Millisecond,
		}
		text.history.max = d.maxPropertyHistory

		d.PropertyText[name] = &text
		d.TextNames = append(d.TextNames, name)
	} else {
		log.Println("overwriting text", name)
		prop.lock.Lock()
		prop.data = *property
		prop.lock.Unlock()
	}
}

func (d *Device) TextValue(name string) (string, bool) {
	property, ok := d.PropertyText[name]
	if ok {
		property.lock.RLock()
		value := property.data.Value
		property.lock.RUnlock()
		return value, true
	}

	return "", false
}

// Was UpdateText
func (d *Device) SetTextValue(name string, value string) {
	fmt.Println("set text", name, value)
	property, ok := d.PropertyText[name]
	if ok {
		property.lock.Lock()
		if property.history.index >= property.history.max {
			property.history.index = 0
		}

		if property.history.max-1 >= len(property.history.values) {
			property.history.values = append(property.history.values, property.data.Value)
		} else {
			property.history.values[property.history.index] = property.data.Value
		}
		property.history.index++

		property.data.Value = value
		// copy the Id so we can unlock before we start the call back action, this means we dont have to
		//  keep the lock open until the client has rwsponded

		property.lock.Unlock()

		if d.clientConnection != nil {
			if writer := d.clientConnection.ClientWriter(d.ClientId); writer != nil {
				jsonOut := d.MakeAction(d.Id, name, DIAL, fmt.Sprint(value))
				writer.Write(jsonOut)
			}
		}
	}

}

func (d *Device) TextWindow(name string, timestamp time.Time) bool {
	property, ok := d.PropertyText[name]
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

func (d *Device) SetTextWindow(name string, duration int64) {
	property, ok := d.PropertyText[name]
	if ok {
		timelimit := time.Duration(duration) * time.Millisecond
		property.lock.Lock()
		property.repeatWindowDuration = timelimit
		property.lock.Unlock()
	}
}

// Was ReadPropertyText
func (d *Device) Map2Text(props map[string]interface{}) (*TextProperty, error) {
	var prop TextProperty
	var err error

	log.Println("reading text property")
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
		prop.Value = v.(string)
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
