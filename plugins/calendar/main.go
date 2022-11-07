package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"
)

type settings struct {
	SockAddr string
}

type Calendar struct {
	c *calendar
	p *plugin
}

const (
	FLAG_NOTSET = iota
	FLAG_HOUR
	FLAG_DAY
	FLAG_WEEK
	FLAG_MONTH
	FLAG_YEAR
)

type calEvent struct {
	msg         string
	repeatCount int
	repeatEvery int
}

func main() {

	profile, err := os.UserConfigDir()
	if err != nil {
		log.Panic("unable to get users home folder", err)
	}
	configPath := path.Join(profile, "auto-home", "plugin.calendar.json")

	jsonFile, err := os.Open(configPath)
	if err != nil {
		log.Println("unable to open plugin.calendar.json", err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var conf settings
	json.Unmarshal(byteValue, &conf)

	p := Connect(SockAddr)

	cal := Calendar{p: p}
	cal.c = New(cal.fireEvent)

	cal.c.Add(time.Now().Add(6*time.Second), "check 3")
	cal.c.Add(time.Now().Add(4*time.Second), "check 2")
	cal.c.Add(time.Now().Add(8*time.Second), "check 4")
	cal.c.Add(time.Now().Add(2*time.Second), "check 1")

	p.Register("calendar", cal)

	err = p.Done()
	if err != nil {
		fmt.Println(err)
	}

}

func (c *Calendar) fireEvent(event interface{}) {
	evt := event.(calendarEvent)
	fmt.Println("fire event:", evt.date, evt.data)

	c.p.Call("onevent", evt)
}
