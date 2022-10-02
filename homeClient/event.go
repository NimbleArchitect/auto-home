package homeClient

import (
	"errors"
	"fmt"
)

type event struct {
	props map[string]property
}

type property struct {
	kind        int
	value       interface{}
	description string
	mode        string
	json        string
}

func NewEvent() event {

	return event{
		props: make(map[string]property),
	}
}

func (e *event) AddDial(name string, value int) error {
	if len(name) == 0 {
		return errors.New("invalid name")
	}

	if _, ok := e.props[name]; ok {
		return errors.New("property exists with that name")
	} else {
		e.props[name] = property{
			kind:  FLAG_DIAL,
			value: value,
			json:  fmt.Sprintf("{\"name\":\"%s\",\"type\":\"dial\",\"value\":%d},", name, value),
		}
		return nil
	}
}

func (e *event) AddSwitch(name string, state interface{}) error {
	if len(name) == 0 {
		return errors.New("invalid name")
	}

	if state == nil {
		return errors.New("state cannot be nil")
	}

	if _, ok := e.props[name]; ok {
		return errors.New("property exists with that name")
	} else {
		e.props[name] = property{
			kind:  FLAG_SWITCH,
			value: state,
			json:  fmt.Sprintf("{\"name\":\"%s\",\"type\":\"switch\",\"value\":\"%s\"},", name, state),
		}
		return nil
	}
}

func (e *event) AddButton(name string, value bool) error {
	if len(name) == 0 {
		return errors.New("invalid name")
	}

	if _, ok := e.props[name]; ok {
		return errors.New("property exists with that name")
	} else {
		e.props[name] = property{
			kind:  FLAG_BUTTON,
			value: value,
			json:  fmt.Sprintf("{\"name\":\"%s\",\"type\":\"button\",\"value\":%t},", name, value),
		}
		return nil
	}
}

func (e *event) AddText(name string, value string) error {
	if len(name) == 0 {
		return errors.New("invalid name")
	}

	if _, ok := e.props[name]; ok {
		return errors.New("property exists with that name")
	} else {
		e.props[name] = property{
			kind:  FLAG_TEXT,
			value: value,
			json:  fmt.Sprintf("{\"name\":\"%s\",\"type\":\"button\",\"value\":%s},", name, value),
		}
		return nil
	}
}
