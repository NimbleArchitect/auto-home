package js

import "log"

type jsSwitch struct {
	Name     string
	Value    bool
	label    string
	previous string
}

func (d *jsSwitch) IsSwitch() bool {
	return true
}
func (d *jsSwitch) IsDial() bool {
	return false
}
func (d *jsSwitch) IsButton() bool {
	return false
}
func (d *jsSwitch) IsText() bool {
	return false
}

func (d *jsSwitch) Type() string {
	return "switch"
}

func (d *jsSwitch) AsBool(name string) bool {
	return d.Value
}

func (d *jsSwitch) Last(x int) interface{} {
	log.Println("TODO: not implemented")
	return nil
}
