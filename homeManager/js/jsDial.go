package js

import "log"

type jsDial struct {
	Name     string
	Value    int
	previous int
	min      int
	max      int
}

func (d *jsDial) IsSwitch() bool {
	return false
}
func (d *jsDial) IsDial() bool {
	return true
}
func (d *jsDial) IsButton() bool {
	return true
}
func (d *jsDial) IsText() bool {
	return true
}

func (d *jsDial) Type() string {
	return "dial"
}

func (d *jsDial) AsPercent() int {

	v := d.Value - d.min
	m := d.max - d.min
	p := (v / m) * 100

	return p
}

func (d *jsDial) Last(x int) interface{} {
	log.Println("TODO: not implemented")
	return nil
}
