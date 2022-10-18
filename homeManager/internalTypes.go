package home

import (
	js "server/homeManager/js"
)

type Hub struct {
	Id          string
	Name        string
	Description string
	ClientId    string
	Help        string
	Devices     []string
}

// type group struct {
// 	Id          string
// 	Name        string
// 	Description string
// 	Devices     []string
// 	Groups      []string
// 	Users       []string
// }

type Action struct {
	Name     string
	Location string
	jsvm     *js.JavascriptVM
}
