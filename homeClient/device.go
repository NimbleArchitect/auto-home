package homeClient

import (
	"errors"
	"fmt"
)

type Device struct {
	Name        string
	Description string
	Id          string
	Help        string
	props       map[string]string
}

// NewDevice returns a new device object, with the name and id fields set
func NewDevice(name string, deviceid string) Device {
	return Device{
		Name:  name,
		Id:    deviceid,
		props: make(map[string]string),
	}
}

func (d *Device) getJson() string {

	var propJson string
	for _, v := range d.props {
		propJson += v
	}
	propJson = propJson[0 : len(propJson)-1]

	jsonData := fmt.Sprintf("{\"id\":\"%s\",\"name\":\"%s\",\"description\":\"%s\",\"help\":\"%s\",\"properties\":[%s]}",
		d.Id, d.Name, d.Description, d.Help, propJson)

	return jsonData
}

func (d *Device) AddDial(name string, description string, value int, min int, max int, mode string) error {
	if len(name) == 0 {
		return errors.New("invalid name")
	}

	if _, ok := d.props[name]; ok {
		return errors.New("property exists with that name")
	} else {
		d.props[name] = fmt.Sprintf("{\"name\":\"%s\",\"description\":\"%s\",\"type\":\"dial\",\"value\":%d,\"min\":%d,\"max\":%d,\"mode\":\"%s\"},", name, description, value, min, max, mode)
		return nil
	}
}

func (d *Device) AddSwitch(name string, description string, state interface{}, mode string) error {
	if len(name) == 0 {
		return errors.New("invalid name")
	}

	if state == nil {
		return errors.New("state cannot be nil")
	}

	if _, ok := d.props[name]; ok {
		return errors.New("property exists with that name")
	} else {
		// TODO: needs fixing converstion to string is lazy
		s := fmt.Sprint(state)
		d.props[name] = fmt.Sprintf("{\"name\":\"%s\",\"description\":\"%s\",\"type\":\"switch\",\"value\":\"%s\",\"mode\":\"%s\"},", name, description, s, mode)
		return nil
	}
}

func (d *Device) AddButton(name string, description string, value bool, normally bool, mode string) error {
	if len(name) == 0 {
		return errors.New("invalid name")
	}

	if _, ok := d.props[name]; ok {
		return errors.New("property exists with that name")
	} else {
		d.props[name] = fmt.Sprintf("{\"name\":\"%s\",\"description\":\"%s\",\"type\":\"button\",\"value\":%t,\"normally\":%t,\"mode\":\"%s\"},", name, description, value, normally, mode)
		return nil
	}
}

func (d *Device) AddText(name string, description string, value string, mode string) error {
	if len(name) == 0 {
		return errors.New("invalid name")
	}

	if _, ok := d.props[name]; ok {
		return errors.New("property exists with that name")
	} else {
		d.props[name] = fmt.Sprintf("{\"name\":\"%s\",\"description\":\"%s\",\"type\":\"text\",\"value\":%s,\"mode\":\"%s\"},", name, description, value, mode)
		return nil
	}
}
