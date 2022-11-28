package main

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"log"
	"os"
	"path"

	"github.com/nathan-osman/go-sunrise"
)

type Client int

type Solar struct {
	conf settings
}

type settings struct {
	Lat float64
	Lon float64
}

// type event struct {
// 	Label    string
// 	Date     time.Time
// 	Location string
// 	Notes    string
// }

func main() {
	profile, err := os.UserConfigDir()
	if err != nil {
		log.Panic("unable to get users home folder", err)
	}
	configPath := path.Join(profile, "auto-home", "plugin.solar.json")

	jsonFile, err := os.Open(configPath)
	if err != nil {
		log.Println("unable to open plugin.solar.json", err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var conf settings
	json.Unmarshal(byteValue, &conf)

	p := Connect()

	cal := new(Solar)
	cal.conf = conf

	if conf.Lat == 0 && conf.Lon == 0 {
		log.Panic("invalid latitude and longitude")
	}

	p.Register("solar", cal)

	go func() {
		currentLight := cal.IsLight()
		for {
			time.Sleep(1 * time.Second)
			newLight := cal.IsLight()
			if currentLight != newLight {
				currentLight = newLight
				if newLight {
					fmt.Println("onSunrise")
					p.Call("onSunrise", nil)
				} else {
					fmt.Println("onSunset")
					p.Call("onSunset", nil)
				}
			}
		}
	}()

	err = p.Done()
	if err != nil {
		fmt.Println(err)
	}

}

func (s *Solar) IsLight() bool {

	now := time.Now()
	rise, set := sunrise.SunriseSunset(s.conf.Lat, s.conf.Lon, now.Year(), now.Month(), now.Day())

	if now.After(rise) && now.Before(set) {
		return true
	}

	return false
}

func (s *Solar) IsDark() bool {
	return !s.IsLight()
}
