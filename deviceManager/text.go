package deviceManager

import (
	"log"
	"sync"
	"time"
)

type Text struct {
	lock sync.RWMutex
	data TextProperty
}
type TextProperty struct {
	Id                    string
	Name                  string
	Description           string
	Value                 string
	Previous              string
	Mode                  uint
	repeatWindowTimeStamp time.Time
	repeatWindowDuration  time.Duration
}

func (d *Device) NewText() *Text {
	return &Text{
		lock: sync.RWMutex{},
		data: TextProperty{},
	}
}

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
		d.PropertyText[name] = &Text{
			lock: sync.RWMutex{},
			data: *property,
		}
	} else {
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
	property, ok := d.PropertyText[name]
	if ok {
		property.lock.Lock()
		property.data.Value = value
		property.lock.Unlock()
	}

}

func (d *Device) TextWindow(name string, timestamp time.Time) bool {
	property, ok := d.PropertyText[name]
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

func (d *Device) SetTextWindow(name string, duration int64) {
	property, ok := d.PropertyText[name]
	if ok {
		property.lock.Lock()
		property.data.repeatWindowDuration = time.Duration(duration) * time.Millisecond
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
