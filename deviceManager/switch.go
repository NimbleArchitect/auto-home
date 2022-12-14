package deviceManager

import (
	"fmt"
	"server/booltype"
	log "server/logger"
	"sync"
	"time"
)

type Switch struct {
	lock    sync.RWMutex
	data    SwitchProperty
	history struct {
		max    int
		index  int
		values []booltype.BoolType
	}
	repeatWindowTimeStamp time.Time
	repeatWindowDuration  time.Duration
}
type SwitchProperty struct {
	// Id          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Value       booltype.BoolType `json:"value"`
	Mode        uint              `json:"mode"`
	Kind        string            `json:"type"`
	// Previous              booltype.BoolType
}

// func (d *Device) NewSwitch() *Switch {
// 	return &Switch{
// 		lock: sync.RWMutex{},
// 		data: SwitchProperty{},
// 	}
// }

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
		duration := d.repeatWindow[name]
		log.Info("new switch", name, duration)

		swi := Switch{
			lock:                 sync.RWMutex{},
			data:                 *property,
			repeatWindowDuration: time.Duration(duration) * time.Millisecond,
		}
		swi.history.max = d.maxPropertyHistory

		d.PropertySwitch[name] = &swi
		d.SwitchNames = append(d.SwitchNames, name)
	} else {
		log.Info("overwriting switch", name)
		prop.lock.Lock()
		prop.data = *property
		prop.lock.Unlock()
	}
}

func (d *Device) SwitchValue(name string) (booltype.BoolType, bool) {
	property, ok := d.PropertySwitch[name]
	if ok {
		property.lock.RLock()
		data := property.data
		property.lock.RUnlock()
		return data.Value, true
	}

	return booltype.BoolType{}, false
}

// updates the live device
func (d *Device) WriteSwitchValue(name string, value string) {

	log.Debug("d.Id =", d.Id)
	if d.clientConnection != nil {
		if writer := d.clientConnection.ClientWriter(d.ClientId); writer != nil {
			jsonOut := d.MakeAction(d.Id, name, SWITCH, fmt.Sprint(value))
			writer.Write(jsonOut)
		}
	} else {
		log.Panic("clientConnection should not be empty")
	}
}

// SetSwitchValue updates the internal value and calls the writer to send the updated value back to the cient
func (d *Device) SetSwitchValue(name string, value string) {

	log.Info("set switch", name, value)

	property, ok := d.PropertySwitch[name]
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

		property.data.Value.Set(value)
		// copy the Id so we can unlock before we start the call back action, this means we dont have to
		//  keep the lock open until the client has rwsponded
		// TODO: is the property id even needed any more??
		// id := property.data.Id

		property.lock.Unlock()
	}
}

func (d *Device) SwitchWindow(name string, timestamp time.Time) bool {
	property, ok := d.PropertySwitch[name]
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

	prop.Kind = "switch"
	log.Info("reading switch property")
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
			log.Error(err)
		}
	}

	return &prop, nil
}
