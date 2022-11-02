package js

import (
	"errors"
	"log"
	"server/booltype"
	"strings"

	"github.com/dop251/goja"
)

type DeviceUpdator interface {
	GetNextVM() (*JavascriptVM, int)
	PushVMID(int)

	// UpdateDial(string, string, int) error
	// UpdateSwitch(string, string, string) error
	// UpdateButton(string, string, string) error
	// UpdateText(string, string, string) error

	// GetDialValue(string, string) (int, bool)
	// GetSwitchValue(string, string) (string, bool)
	// GetButtonValue(string, string) (string, bool)
	// GetTextValue(string, string) (string, bool)

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
func (r *JavascriptVM) runAsThread(function goja.Value, value goja.Value) {
	go func() {
		var jsHome jsHome
		var ok bool

		vm, id := r.Updater.GetNextVM()
		defer r.Updater.PushVMID(id)

		jsHome.StopProcessing = FLAG_STOPPROCESSING
		jsHome.ContinueProcessing = FLAG_CONTINUEPROCESSING
		jsHome.GroupProcessing = FLAG_GROUPPROCESSING

		jsHome.devices = r.deviceState

		vm.runtime.Set("plugin", vm.plugins)
		vm.runtime.Set("home", jsHome)

		obj := vm.runtime.Get(function.String())
		call, ok := goja.AssertFunction(obj)
		if ok {
			if value == nil {
				call(goja.Undefined())
			} else {
				call(goja.Undefined(), value)
			}
		} else {
			log.Println("thread call not a function")
		}
	}()
}

func MapToJsSwitch(prop map[string]interface{}) (jsSwitch, error) {
	var swi jsSwitch
	var tmpBool booltype.BoolType

	if n, ok := prop["name"]; ok {
		swi.Name = n.(string)
	} else {
		return jsSwitch{}, errors.New("missing name")
	}

	if b, ok := prop["value"]; ok {
		tmpBool.Set(b.(string))
		swi.Value = tmpBool.String()
		swi.state = tmpBool.GetBool()
	} else {
		return jsSwitch{}, errors.New("missing value")
	}

	return swi, nil
}

func MapToJsDial(prop map[string]interface{}) (jsDial, error) {
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
func MapToJsButton(prop map[string]interface{}) (jsButton, error) {
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

func MapToJsText(prop map[string]interface{}) (jsText, error) {
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
