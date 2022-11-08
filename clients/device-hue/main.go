package main

import (
	"crypto/tls"
	homeClient "device-hue/homeClient"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"golang.org/x/net/http2"
)

type settings struct {
	Username   string
	HubAddress string
	devices    map[string]string
	http       *http.Client
	ServerURL  string
}

func main() {
	const token = "randomhuehubuuid"

	profile, err := os.UserConfigDir()
	if err != nil {
		log.Panic("unable to get users home folder", err)
	}
	configPath := path.Join(profile, "auto-home", "device.hue.json")

	jsonFile, err := os.Open(configPath)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var conf settings
	json.Unmarshal(byteValue, &conf)

	client := homeClient.NewClient(conf.ServerURL, token)

	transport := &http2.Transport{
		TLSClientConfig: &tls.Config{
			// RootCAs:            pool,
			InsecureSkipVerify: true,
			// KeyLogWriter:       keyLog,
		},
		DisableCompression: true,
		AllowHTTP:          false,
	}

	conf.http = &http.Client{Transport: transport}

	conf.hueRegisterHub(conf.Username, conf.HubAddress, &client)

	event, err := client.ListenEvents(conf.callback)
	if err != nil {
		log.Panic("unable to listen", err)

	}

	finished := false
	for {
		select {
		case msg := <-event:
			switch msg {
			case homeClient.EVENT_RELOAD:
				conf.hueRegisterHub(conf.Username, conf.HubAddress, &client)

			case homeClient.EVENT_SHUTDOWN:
				finished = true
			default:
				fmt.Println(">> *** hit default ***")
			}
		case <-time.After(10 * time.Second):
			log.Println(">> pull state <<")
		}
		if finished {
			break
		}
	}

	log.Println("!>> got here")
	// evt := homeClient.NewEvent()

	// evt.AddDial("hue", 50)
	// evt.AddSwitch("state", "on")

	// client.SendEvent("123-tv-light-321", evt)

	// time.Sleep(600 * time.Second)
}

func (s *settings) callback(deviceid string, args map[string]interface{}) {
	var j string

	// fmt.Println(">>", s.HubAddress)

	id := s.devices[deviceid]
	url := s.HubAddress + "/api/" + s.Username + "/lights/"

	for k, v := range args {
		// fmt.Println("!!>>", k, v)
		// s.http
		if k == "brightness" {
			f, _ := v.(float64)
			if f <= 0 {
				j = `{"on":false}`
			} else {
				s := strconv.FormatFloat(f, 'f', 0, 32)
				j = `{"on":true,"bri":` + s + `}`
			}
		}
		if k == "state" {
			s := v.(string)
			j = `{"on":` + s + `}`
		}

		// fmt.Println(url + id + "/state")
		// fmt.Println(">> posting <<", j)
		_, err := s.Put(url+id+"/state", j)
		if err != nil {
			log.Println("error sending /state", err)
		}
	}
}
