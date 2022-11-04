package js

import "log"

type jsButton struct {
	Name     string
	Value    string //new .state, old .label
	state    bool   //old .Value
	previous string
	flag     jsFlag
}

func (d *jsButton) IsSwitch() bool {
	return false
}

func (d *jsButton) IsDial() bool {
	return false
}

func (d *jsButton) IsButton() bool {
	return true
}

func (d *jsButton) IsText() bool {
	return false
}

func (d *jsButton) Type() string {
	return "button"
}

func (d *jsButton) AsBool(name string) bool {
	return d.state
}

func (d *jsButton) Last(x int) interface{} {
	log.Println("TODO: not implemented")
	return nil
}
