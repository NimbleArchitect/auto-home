package js

import (
	"errors"
	"log"
	"server/booltype"
	"strings"

	"github.com/dop251/goja"
)

type DeviceUpdator interface {
	UpdateDial(string, string, int) error
	UpdateSwitch(string, string, string) error
	UpdateButton(string, string, string) error
	UpdateText(string, string, string) error

	GetDialValue(string, string) (int, bool)
	GetSwitchValue(string, string) (string, bool)
	GetButtonValue(string, string) (string, bool)
	GetTextValue(string, string) (string, bool)

	// GetDialHistory()

	// RunGroupAction(string, string, []map[string]interface{}) (interface{}, error)
}

const (
	StrOnTrigger = "ontrigger"
	StrOnChange  = "onchange"
	StrOnStart   = "onstart"
)

func BuildOnAction(values ...string) string {
	if len(values) == 1 {
		return values[0]
	}

	return strings.Join(values, "_")
}

// runAsThread runs the js function as a new thread, this could be dangerous/not thread safe
func runAsThread(obj goja.Value, val goja.Value) {
	call, ok := goja.AssertFunction(obj)
	if ok {
		go call(goja.Undefined(), val)
	} else {
		log.Println("thread call not a function")
	}
}

func mapToJsSwitch(prop map[string]interface{}) (jsSwitch, error) {
	var swi jsSwitch
	var tmpBool booltype.BoolType

	if n, ok := prop["name"]; ok {
		swi.Name = n.(string)
	} else {
		return jsSwitch{}, errors.New("missing name")
	}

	if b, ok := prop["value"]; ok {
		tmpBool.Set(b.(string))
		swi.label = tmpBool.String()
		swi.Value = tmpBool.GetBool()
	} else {
		return jsSwitch{}, errors.New("missing value")
	}

	return swi, nil
}

func mapToJsDial(prop map[string]interface{}) (jsDial, error) {
	var dial jsDial

	if n, ok := prop["name"]; ok {
		dial.Name = n.(string)
	} else {
		return jsDial{}, errors.New("missing name")
	}

	if c, ok := prop["value"]; ok {
		f, isFloat := c.(float64)
		if !isFloat {
			return jsDial{}, errors.New("error converting property current")
		}
		dial.Value = int(f)
	} else {
		return jsDial{}, errors.New("missing value")
	}

	return dial, nil
}

// TODO: add button and text props
func mapToJsButton(prop map[string]interface{}) (jsButton, error) {
	var button jsButton

	if n, ok := prop["name"]; ok {
		button.Name = n.(string)
	} else {
		return jsButton{}, errors.New("missing name")
	}

	if c, ok := prop["value"]; ok {
		f, isBool := c.(bool)
		if !isBool {
			return jsButton{}, errors.New("error converting property current")
		}
		button.Value = f
	} else {
		return jsButton{}, errors.New("missing value")
	}

	return button, nil
}

func mapToJsText(prop map[string]interface{}) (jsText, error) {
	var text jsText

	if n, ok := prop["name"]; ok {
		text.Name = n.(string)
	} else {
		return jsText{}, errors.New("missing name")
	}

	if c, ok := prop["value"]; ok {
		text.Value = c.(string)
	} else {
		return jsText{}, errors.New("missing value")
	}

	return text, nil
}
