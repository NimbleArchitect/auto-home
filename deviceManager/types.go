package deviceManager

import "errors"

const (
	BUTTON = iota
	DIAL
	SWITCH
	TEXT

	RO
	RW
	WO
)

// type group struct {
// 	Id          string
// 	Name        string
// 	Description string
// 	Devices     []string
// 	Groups      []string
// 	Users       []string
// }

// type Upload struct {
// 	Name  string
// 	Alias []string
// }

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

type ActionWriter interface {
	Write(s string) (int, error)
}
