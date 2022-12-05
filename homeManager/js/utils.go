package js

import (
	"errors"

	"server/booltype"
	log "server/logger"
	"strconv"
	"strings"

	"github.com/dop251/goja"
)

const (
	StrOnTrigger = "ontrigger"
	StrOnChange  = "onchange"
	StrOnStart   = "onstart"
)

type DeviceUpdator interface {
	GetNextVM() (*JavascriptVM, int)
	PushVMID(int)
}

func BuildOnAction(values ...string) string {
	if len(values) == 1 {
		return values[0]
	}

	return strings.Join(values, "_")
}

// runAsThread runs the js function as a new thread, this could be dangerous/not thread safe
func (r *JavascriptVM) runAsThread(function goja.Value, value goja.Value) {
	r.vmInUseLock.Add(1)

	go func() {
		var ok bool

		vm, id := r.Updater.GetNextVM()
		defer r.Updater.PushVMID(id)

		vm.setJsGlobal()

		obj := vm.runtime.Get(function.String())
		call, ok := goja.AssertFunction(obj)
		if ok {
			if value == nil {
				call(goja.Undefined())
			} else {
				call(goja.Undefined(), value)
			}
		} else {
			log.Error("thread call not a function")
		}
		r.vmInUseLock.Done()
	}()
}

func MapToJsSwitch(prop map[string]interface{}) (jsSwitch, error) {
	var swi jsSwitch
	var tmpBool booltype.BoolType

	if n, ok := prop["name"]; ok {
		swi.Name, ok = n.(string)
		if !ok {
			return jsSwitch{}, errors.New("name is not a string")
		}
	} else {
		return jsSwitch{}, errors.New("missing name")
	}

	if c, ok := prop["value"]; ok {
		f, isBool := c.(bool)
		if isBool {
			tmpBool.SetBool(f)
			swi.state = tmpBool.Bool()
			swi.Value = tmpBool.String()
		} else {
			s, isString := c.(string)
			if isString {
				tmpBool.Set(s)
				swi.state = tmpBool.Bool()
				swi.Value = tmpBool.String()
			} else {
				return jsSwitch{}, errors.New("error converting property value")
			}
		}
	} else {
		return jsSwitch{}, errors.New("missing value")
	}

	return swi, nil
}

func MapToJsDial(prop map[string]interface{}) (jsDial, error) {
	var dial jsDial

	if n, ok := prop["name"]; ok {
		dial.Name, ok = n.(string)
		if !ok {
			return jsDial{}, errors.New("name is not a string")
		}
	} else {
		return jsDial{}, errors.New("missing name")
	}

	if c, ok := prop["value"]; ok {
		switch value := c.(type) {
		case string:
			i, err := strconv.ParseInt(value, 0, 0)
			if err != nil {
				return jsDial{}, errors.New("error converting property value: " + err.Error())
			}
			dial.Value = int(i)
		case int:
			dial.Value = value
		case float64:
			dial.Value = int(value)
		default:
			return jsDial{}, errors.New("error converting property value")
		}

	} else {
		return jsDial{}, errors.New("missing value")
	}

	return dial, nil
}

// TODO: add button and text props
func MapToJsButton(prop map[string]interface{}) (jsButton, error) {
	var button jsButton
	var tmpBool booltype.BoolType

	if n, ok := prop["name"]; ok {
		button.Name, ok = n.(string)
		if !ok {
			return jsButton{}, errors.New("name is not a string")
		}
	} else {
		return jsButton{}, errors.New("missing name")
	}

	if c, ok := prop["value"]; ok {
		f, isBool := c.(bool)
		if isBool {
			tmpBool.SetBool(f)
			button.state = tmpBool.Bool()
			button.Value = tmpBool.String()
		} else {
			s, isString := c.(string)
			if isString {
				tmpBool.Set(s)
				button.state = tmpBool.Bool()
				button.Value = tmpBool.String()
			} else {
				return jsButton{}, errors.New("error converting property value")
			}
		}

	} else {
		return jsButton{}, errors.New("missing value")
	}

	return button, nil
}

func MapToJsText(prop map[string]interface{}) (jsText, error) {
	var text jsText

	if n, ok := prop["name"]; ok {
		text.Name, ok = n.(string)
		if !ok {
			return jsText{}, errors.New("name is not a string")
		}
	} else {
		return jsText{}, errors.New("missing name")
	}

	if c, ok := prop["value"]; ok {
		t, ok := c.(string)
		if !ok {
			return jsText{}, errors.New("error converting property value")
		}
		text.Value = t
	} else {
		return jsText{}, errors.New("missing value")
	}

	return text, nil
}
