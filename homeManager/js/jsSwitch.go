package js

import (
	"log"
)

type jsSwitch struct {
	Name     string
	Value    string
	state    bool
	previous string
	flag     jsFlag
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

func (d *jsSwitch) ValueOf() bool {
	return d.state
}

func (d *jsSwitch) ToString() string {
	return d.Value
}

func (d *jsSwitch) Last(x int) interface{} {
	log.Println("TODO: not implemented")
	return nil
}
