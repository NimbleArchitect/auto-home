package deviceManager

import (
	"log"
	"server/booltype"
	"sync"
	"time"
)

type Switch struct {
	lock                  sync.RWMutex
	data                  SwitchProperty
	repeatWindowTimeStamp time.Time
	repeatWindowDuration  time.Duration
}
type SwitchProperty struct {
	Id          string
	Name        string
	Description string
	Value       booltype.BoolType
	// Previous              booltype.BoolType
	Mode uint
}

func (d *Device) NewSwitch() *Switch {
	return &Switch{
		lock: sync.RWMutex{},
		data: SwitchProperty{},
	}
}

func (d *Device) SwitchKeys() []string {

	return d.SwitchNames
}

func (d *Device) Switch(name string) SwitchProperty {
	property, ok := d.PropertySwitch[name]
	if ok {
		property.lock.RLock()
		out := property.data
		property.lock.RUnlock()
		return out
	}

	return SwitchProperty{}
}

func (d *Device) SwitchAsMap() map[string]SwitchProperty {
	propertyList := make(map[string]SwitchProperty, len(d.PropertySwitch))

	for name, property := range d.PropertySwitch {
		property.lock.RLock()
		propertyList[name] = property.data
		property.lock.RUnlock()
	}

	return propertyList
}

func (d *Device) SetSwitch(name string, property *SwitchProperty) {
	prop, ok := d.PropertySwitch[name]
	if !ok {
		d.PropertySwitch[name] = &Switch{
			lock: sync.RWMutex{},
			data: *property,
		}
		d.SwitchNames = append(d.SwitchNames, name)
	} else {
		prop.lock.Lock()
		prop.data = *property
		prop.lock.Unlock()
	}
}

func (d *Device) SwitchValue(name string) (string, bool) {
	property, ok := d.PropertySwitch[name]
	if ok {
		property.lock.RLock()
		data := property.data
		property.lock.RUnlock()
		return data.Value.String(), true
	}

	return "", false
}

// Was UpdateSwitch
func (d *Device) SetSwitchValue(name string, value string) {
	if d.PropertySwitch == nil {
		return
	}

	property, ok := d.PropertySwitch[name]
	if ok {
		property.lock.Lock()
		property.data.Value.Set(value)
		property.lock.Unlock()
	}

}

func (d *Device) SwitchWindow(name string, timestamp time.Time) bool {
	property, ok := d.PropertySwitch[name]
	if ok {
		if property.repeatWindowTimeStamp.Before(timestamp) {
			newExpire := timestamp.Add(time.Duration(property.repeatWindowDuration) * time.Millisecond)
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

func (d *Device) SetSwitchWindow(name string, duration int64) {
	property, ok := d.PropertySwitch[name]
	if ok {
		timelimit := time.Duration(duration) * time.Millisecond
		property.lock.Lock()
		property.repeatWindowDuration = timelimit
		property.lock.Unlock()
	}
}

// Was ReadPropertySwitch
func (d *Device) Map2Switch(props map[string]interface{}) (*SwitchProperty, error) {
	var prop SwitchProperty
	var err error

	log.Println("reading switch property")
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
