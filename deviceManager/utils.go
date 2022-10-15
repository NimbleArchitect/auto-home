package deviceManager

import (
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
