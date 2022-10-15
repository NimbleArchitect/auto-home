package deviceManager

import (
	"log"
	"sync"
	"time"
)

type Dial struct {
	lock                  sync.RWMutex
	data                  DialProperty
	repeatWindowTimeStamp time.Time
	repeatWindowDuration  time.Duration
}
type DialProperty struct {
	Id          string
	Name        string
	Description string
	Min         int
	Max         int
	Value       int
	// Previous              int
	Mode uint
}

func (d *Device) NewDial() *Dial {
	return &Dial{
		lock: sync.RWMutex{},
		data: DialProperty{},
	}
}

func (d *Device) DialKeys() []string {

	return d.DialNames
}

func (d *Device) Dial(name string) DialProperty {
	property, ok := d.PropertyDial[name]
	if ok {
		property.lock.RLock()
		out := property.data
		property.lock.RUnlock()
		return out
	}

	return DialProperty{}
}

func (d *Device) DialAsMap() map[string]DialProperty {
	propertyList := make(map[string]DialProperty, len(d.PropertyDial))

	for name, property := range d.PropertyDial {
		property.lock.RLock()
		propertyList[name] = property.data
		property.lock.RUnlock()
	}

	return propertyList
}

func (d *Device) SetDial(name string, property *DialProperty) {
	prop, ok := d.PropertyDial[name]
	if !ok {
		duration := d.repeatWindow[name]
		log.Println("new dial", name, duration)

		d.PropertyDial[name] = &Dial{
			lock:                 sync.RWMutex{},
			data:                 *property,
			repeatWindowDuration: time.Duration(duration) * time.Millisecond,
		}
		d.DialNames = append(d.DialNames, name)
	} else {
		log.Println("overwriting dial", name)
		prop.lock.Lock()
		prop.data = *property
		prop.lock.Unlock()
	}
}

func (d *Device) DialValue(name string) (int, bool) {
	property, ok := d.PropertyDial[name]
	if ok {
		property.lock.RLock()
		data := property.data
		property.lock.RUnlock()
		return data.Value, ok
	}

	return 0, false
}

// Was UpdateDial
func (d *Device) SetDialValue(name string, value int) {
	property, ok := d.PropertyDial[name]
	if ok {
		property.lock.Lock()
		property.data.Value = value
		property.lock.Unlock()
	}

}

func (d *Device) DialWindow(name string, timestamp time.Time) bool {
	property, ok := d.PropertyDial[name]
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

func (d *Device) SetDialWindow(name string, duration int64) {
	property, ok := d.PropertyDial[name]
	if ok {
		timelimit := time.Duration(duration) * time.Millisecond
		property.lock.Lock()
		property.repeatWindowDuration = timelimit
		property.lock.Unlock()
	}
}

// Was ReadPropertyDial
func (d *Device) Map2Dial(props map[string]interface{}) (*DialProperty, error) {
	var prop DialProperty
	var err error

	log.Println("reading dial property")
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

	if v, ok := props["min"]; !ok {
		return nil, ErrMissingPropertyMin
	} else {
		f, isFloat := v.(float64)
		if !isFloat {
			return nil, ErrConvertingPropteryMin
		}
		prop.Min = int(f)
	}

	if v, ok := props["max"]; !ok {
		return nil, ErrMissingPropertyMax
	} else {
		f, isFloat := v.(float64)
		if !isFloat {
			return nil, ErrConvertingPropteryMax
		}
		prop.Max = int(f)
	}

	// if min is bigger than max swap them around
	if prop.Max < prop.Min {
		tmp := prop.Max
		prop.Max = prop.Min
		prop.Max = tmp
	}

	if v, ok := props["value"]; !ok {
		return nil, ErrMissingPropertyValue
	} else {
		f, isFloat := v.(float64)
		if !isFloat {
			return nil, ErrConvertingPropteryValue
		}
		prop.Value = int(f)
		if prop.Value > prop.Max {
			prop.Value = prop.Max
		}
		if prop.Value < prop.Min {
			prop.Value = prop.Min
		}
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
