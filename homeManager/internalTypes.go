package home

import (
	"server/booltype"
	js "server/homeManager/js"
)

const (
	BUTTON = iota
	DIAL
	SWITCH
	TEXT

	RO
	RW
	WO
)

type Hub struct {
	Id          string
	Name        string
	Description string
	ClientId    string
	Help        string
	Devices     []string
}

type Device struct {
	Id             string
	Name           string
	Description    string
	ClientId       string
	Help           string
	Groups         []group
	PropertySwitch map[string]SwitchProperty
	PropertyDial   map[string]DialProperty
	Uploads        []Upload
}

type group struct {
	Id          string
	Name        string
	Description string
	Device      []string
	Actions     map[string]string
}

type Upload struct {
	Name  string
	Alias []string
}

type Action struct {
	Name     string
	Location string
	jsvm     *js.JavascriptVM
}

type DialProperty struct {
	Name        string
	Description string
	Min         int
	Max         int
	Value       int
	Previous    int
	Mode        uint
}

type SwitchProperty struct {
	Name        string
	Description string
	Value       booltype.BoolType
	Previous    booltype.BoolType
	Mode        uint
}
