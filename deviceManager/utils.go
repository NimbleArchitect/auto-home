package deviceManager

import (
	"errors"
	"strings"
)

var (
	ErrMissingPropertyName  = errors.New("missing property name")
	ErrMissingPropertyValue = errors.New("missing property value")
	ErrMissingPropertyMode  = errors.New("missing property mode")
	ErrMissingPropertyMin   = errors.New("missing property min")
	ErrMissingPropertyMax   = errors.New("missing property max")

	ErrConvertingPropteryMin   = errors.New("error converting property min")
	ErrConvertingPropteryMax   = errors.New("error converting property max")
	ErrConvertingPropteryValue = errors.New("error converting property value")

	ErrInvalidModeValue = errors.New("invalid value for property mode")

	ErrWriteOnlyProperty = errors.New("unable to read from write only property")
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
