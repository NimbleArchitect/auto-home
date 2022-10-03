package js

import (
	"errors"
	"log"

	booltype "server/booltype"

	"github.com/dop251/goja"
)

func (r *JavascriptVM) RunJS(fName string, props goja.Value) (goja.Value, error) {
	var jsHome jsHome

	jsFunction := r.runtime.Get(fName)
	call, ok := goja.AssertFunction(jsFunction)
	if !ok {
		// slient ignore as the function dosent exist in javascript
		log.Println("function", fName, "doesn't exist, skipping")
		return nil, nil
	}

	jsHome.StopProcessing = FLAG_STOPPROCESSING
	jsHome.ContinueProcessing = FLAG_CONTINUEPROCESSING
	jsHome.GroupProcessing = FLAG_GROUPPROCESSING
	jsHome.devices = r.deviceState

	r.runtime.Set("home", jsHome)

	result, err := call(goja.Undefined(), props)
	if err != nil {
		log.Println(err)
	}

	return result, err
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
