package home

import (
	"log"
	"strings"
)

func GetModeFromString(value string) (uint, error) {
	// TODO: have I finished coding RW permissions?
	b := strings.ToLower(value)
	switch b {
	case "ro":
		return RO, nil
	case "rw":
		return RW, nil
	case "wo":
		return WO, nil

	default:
		return RO, ErrInvalidModeValue
	}
}

func ReadPropertySwitch(props map[string]interface{}) (SwitchProperty, error) {
	var prop SwitchProperty
	var err error

	log.Println("reading switch property")
	if v, ok := props["name"]; !ok {
		return SwitchProperty{}, ErrMissingPropertyName
	} else {
		// TODO: clean the string
		prop.Name = v.(string)
		log.Println("name", prop.Name)
	}

	if v, ok := props["description"]; ok {
		// TODO: clean the string
		prop.Description = v.(string)
	}

	if v, ok := props["value"]; !ok {
		return SwitchProperty{}, ErrMissingPropertyValue
	} else {
		prop.Value.Set(v.(string))
	}

	if v, ok := props["mode"]; !ok {
		return SwitchProperty{}, ErrMissingPropertyMode
	} else {
		prop.Mode, err = GetModeFromString(v.(string))
		if err != nil {
			log.Println(err)
		}
	}

	return prop, nil
}

func ReadPropertyDial(props map[string]interface{}) (DialProperty, error) {
	var prop DialProperty
	var err error

	log.Println("reading dial property")
	if v, ok := props["name"]; !ok {
		return DialProperty{}, ErrMissingPropertyName
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
		return DialProperty{}, ErrMissingPropertyMin
	} else {
		f, isFloat := v.(float64)
		if !isFloat {
			return DialProperty{}, ErrConvertingPropteryMin
		}
		prop.Min = int(f)
	}

	if v, ok := props["max"]; !ok {
		return DialProperty{}, ErrMissingPropertyMax
	} else {
		f, isFloat := v.(float64)
		if !isFloat {
			return DialProperty{}, ErrConvertingPropteryMax
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
		return DialProperty{}, ErrMissingPropertyValue
	} else {
		f, isFloat := v.(float64)
		if !isFloat {
			return DialProperty{}, ErrConvertingPropteryValue
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
		return DialProperty{}, ErrMissingPropertyMode
	} else {
		prop.Mode, err = GetModeFromString(v.(string))
		if err != nil {
			log.Println(err)
		}
	}

	return prop, nil
}

func ReadPropertyButton(props map[string]interface{}) (ButtonProperty, error) {
	var prop ButtonProperty
	var err error

	log.Println("reading button property")
	if v, ok := props["name"]; !ok {
		return ButtonProperty{}, ErrMissingPropertyName
	} else {
		// TODO: clean the string
		prop.Name = v.(string)
		log.Println("name", prop.Name)
	}

	if v, ok := props["description"]; ok {
		// TODO: clean the string
		prop.Description = v.(string)
	}

	if v, ok := props["value"]; !ok {
		return ButtonProperty{}, ErrMissingPropertyValue
	} else {
		prop.Value.Set(v.(string))
	}

	if v, ok := props["mode"]; !ok {
		return ButtonProperty{}, ErrMissingPropertyMode
	} else {
		prop.Mode, err = GetModeFromString(v.(string))
		if err != nil {
			log.Println(err)
		}
	}

	return prop, nil
}

func ReadPropertyText(props map[string]interface{}) (TextProperty, error) {
	var prop TextProperty
	var err error

	log.Println("reading text property")
	if v, ok := props["name"]; !ok {
		return TextProperty{}, ErrMissingPropertyName
	} else {
		// TODO: clean the string
		prop.Name = v.(string)
		log.Println("name", prop.Name)
	}

	if v, ok := props["description"]; ok {
		// TODO: clean the string
		prop.Description = v.(string)
	}

	if v, ok := props["value"]; !ok {
		return TextProperty{}, ErrMissingPropertyValue
	} else {
		prop.Value = v.(string)
	}

	if v, ok := props["mode"]; !ok {
		return TextProperty{}, ErrMissingPropertyMode
	} else {
		prop.Mode, err = GetModeFromString(v.(string))
		if err != nil {
			log.Println(err)
		}
	}

	return prop, nil
}
