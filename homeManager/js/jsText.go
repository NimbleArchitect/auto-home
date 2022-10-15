package js

import "log"

type jsText struct {
	Name     string
	Value    string
	previous string
}

func (d *jsText) IsSwitch() bool {
	return false
}
func (d *jsText) IsDial() bool {
	return false
}
func (d *jsText) IsButton() bool {
	return false
}
func (d *jsText) IsText() bool {
	return true
}

func (d *jsText) Type() string {
	return "text"
}

func (d *jsText) AsBool(name string) string {
	return d.Value
}

func (d *jsText) Last(x int) interface{} {
	log.Println("TODO: not implemented")
	return nil
}
