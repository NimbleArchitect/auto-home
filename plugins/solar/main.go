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

// const SockAddr = "/tmp/rpc.sock"

type Client int

type Solar struct {
	conf settings
}

type settings struct {
	SockAddr string
	Lat      float64
	Lon      float64
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

	p := Connect(SockAddr)

	cal := new(Solar)
	cal.conf = conf
	p.Register("solar", cal)

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
