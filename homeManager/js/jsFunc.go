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
	jsTrigger, ok := goja.AssertFunction(jsFunction)
	if !ok {
		// slient ignore as the function dosent exist in javascript
		log.Println("function", fName, "doesn't exist, skipping")
		return nil, nil
	}

	jsHome.devices = r.deviceState

	r.runtime.Set("home", jsHome)

	result, err := jsTrigger(goja.Undefined(), props)
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
