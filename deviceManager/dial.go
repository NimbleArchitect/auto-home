package deviceManager

import (
	"fmt"
	"server/logger"
	"sync"
	"time"
)

type Dial struct {
	lock    sync.RWMutex
	data    DialProperty
	history struct {
		max    int
		index  int
		values []int
	}
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

// func (d *Device) NewDial(maxHistory int) *Dial {
// 	return &Dial{
// 		lock: sync.RWMutex{},
// 		data: DialProperty{},
// 	}
// }

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
	log := logger.New("SetDial", &debugLevel)

	prop, ok := d.PropertyDial[name]
	if !ok {
		duration := d.repeatWindow[name]
		log.Info("new dial", name, duration)

		dial := Dial{
			lock:                 sync.RWMutex{},
			data:                 *property,
			repeatWindowDuration: time.Duration(duration) * time.Millisecond,
		}
		dial.history.max = d.maxPropertyHistory

		d.PropertyDial[name] = &dial
		d.DialNames = append(d.DialNames, name)
	} else {
		log.Info("overwriting dial", name)
		prop.lock.Lock()
		prop.data = *property
		prop.lock.Unlock()
	}
}

func (d *Device) DialValue(name string) (int, bool) {
	property, ok := d.PropertyDial[name]
	if ok {
		property.lock.RLock()
		value := property.data.Value
		property.lock.RUnlock()
		return value, ok
	}

	return 0, false
}

// updates the live device
func (d *Device) WriteDialValue(name string, value int) {
	log := logger.New("WriteDialValue", &debugLevel)
	log.Debug("d.Id =", d.Id)

	if d.clientConnection != nil {
		if writer := d.clientConnection.ClientWriter(d.ClientId); writer != nil {
			jsonOut := d.MakeAction(d.Id, name, DIAL, fmt.Sprint(value))
			writer.Write(jsonOut)
		}
	}
}

// Was UpdateDial
func (d *Device) SetDialValue(name string, value int) {
	log := logger.New("SetDialValue", &debugLevel)
	log.Info("set dial", name, value)

	property, ok := d.PropertyDial[name]
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

	log := logger.New("Map2Dial", &debugLevel)

	log.Info("reading dial property")
	if v, ok := props["name"]; !ok {
		return nil, ErrMissingPropertyName
	} else {
		// TODO: clean the string
		prop.Name = v.(string)
		log.Info("name", prop.Name)
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
			log.Error(err)
		}
	}

	return &prop, nil
}
