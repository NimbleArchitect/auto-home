package home

import (
	js "server/homeManager/js"
	"time"
)

type Hub struct {
	Id          string
	Name        string
	Description string
	ClientId    string
	Help        string
	Devices     []string
}

type group struct {
	Id          string
	Name        string
	Description string
	Devices     []string
	Groups      []string
	Users       []string
}

// type Upload struct {
// 	Name  string
// 	Alias []string
// }

type Action struct {
	Name     string
	Location string
	jsvm     *js.JavascriptVM
}

// notInTimeWindow - checks the suppilied time against the properties update window,
//
//	returns true if we are outside of the update window
//	        false all other times
func notInTimeWindow(time time.Time, propTimeStamp time.Time) bool {
	if time.Before(propTimeStamp) {
		return true
	}
	return false
}
