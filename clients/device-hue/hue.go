package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	homeClient "device-hue/homeClient"

	"golang.org/x/net/http2"
)

type hueInfo struct {
	Lights map[string]light
	Config hueConfig
}

type lightState struct {
	On        bool
	Bri       int
	Alert     string
	Mode      string
	Reachable bool
}

type light struct {
	Name         string
	Type         string
	Modelid      string
	Productname  string
	Productid    string
	Uniqueid     string
	State        lightState
	Capabilities map[string]interface{}
	id           string
}

type hueConfig struct {
	Name     string
	Bridgeid string
}

func (s *settings) hueRegisterHub(username string, url string, client *homeClient.AhClient) {
	info := s.getHueInfo()

	s.State = make(map[string]lightState)

	hub := homeClient.NewHub(info.Config.Name, info.Config.Bridgeid)

	if len(s.devices) == 0 {
		s.devices = make(map[string]light)
	}

	//TODO: this dosent work device properties dont do anything
	for key, v := range info.Lights {
		dev := homeClient.NewDevice(v.Name, v.Uniqueid)

		dev.AddSwitch("state", v.Productname, v.State.On, "RW")
		dev.AddDial("brightness", v.Productname, v.State.Bri, 0, 256, "RW")
		hub.AddDevice(dev)

		v.id = key
		s.devices[v.Uniqueid] = v
	}

	client.RegisterHub(&hub)
}

func (s *settings) getHueInfo() hueInfo {

	res, err := s.http.Get(s.HubAddress + "/api/" + s.Username)
	if err != nil {
		log.Fatal(err)
	}

	var info hueInfo
	err = json.NewDecoder(res.Body).Decode(&info)
	if err != nil {
		log.Println(err)
	}

	defer res.Body.Close()

	return info
}

func (s *settings) Put(url string, data string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return nil, err
	}

	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := s.http.Do(req)
	if err != nil {
		fmt.Println("Put error - url " + url + ": " + err.Error())
		return nil, err
	}

	return resp, nil
}

func (s *settings) listenEvents() {
	// TODO: listen to hue event stream and return items over the channel
	var data string

	transp := http2.Transport{
		ReadIdleTimeout: 20 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	s.http.Transport = &transp

	for {
		isConnected := false
		// to get a list of resources (lights, switches, temps, groups, scenes, dynamic scenes)
		// use /clip/v2/resource
		req, err := http.NewRequest(http.MethodGet, s.HubAddress+"/eventstream/clip/v2", bytes.NewBuffer([]byte{}))
		if err != nil {
			fmt.Println("eventstream read error:", err)
			break
		}
		req.Header.Set("hue-application-key", s.Username)
		req.Header.Set("Accept", "text/event-stream")
		res, err := s.http.Do(req)
		if err != nil {
			log.Fatal("listenEvents() error:", err)
		}

		scanner := bufio.NewScanner(res.Body)
		if scanner.Err() != nil {
			log.Println("unable to start scanner:", scanner.Err())
			break
		}

		for scanner.Scan() {
			// on a read error this loop breaks out

			ln := scanner.Text()
			if len(ln) == 0 {
				//recieved empty line from the server
				// treat it as a connection test and ignore
				fmt.Println("scanner empty")
				continue
			}

			fmt.Println("scanner recieved", ln)

			switch {
			case strings.Contains(ln, ": hi"):
				isConnected = true
				continue
			case strings.HasPrefix(ln, "id: "):
				// id = ln
			case strings.HasPrefix(ln, "data: "):
				data = ln
				if isConnected {
					s.eventDecode(data)
				}
			}

		}
	}
}

type hueEvent struct {
	Creationtime string
	Data         []hueData
	Id           string
	Kind         string `json:"type"`
}

type hueData struct {
	On      map[string]bool
	Dimming map[string]float64
	Palette *struct {
		Color            []hueColor
		Color_temprature []hueColorTemperature
	}
	Color *hueColor
	Id    string
	Id_v1 string
	Owner map[string]string
	Kind  string `json:"type"`
}

type hueColor struct {
	Xy colorFloat
}

type colorFloat struct {
	X float64
	Y float64
}

type hueColorTemperature struct {
	Color_temperature map[string]int
}

type colorTemperature struct {
	Mirek       *int
	Mirek_valid bool
}

func (s *settings) eventDecode(data string) {
	var eventarray []hueEvent

	err := json.Unmarshal([]byte(data[6:]), &eventarray)
	if err != nil {
		fmt.Println("json decode error:", err)
	}

	for _, eventMsg := range eventarray {
		if eventMsg.Kind == "update" {
			for _, event := range eventMsg.Data {
				fmt.Println("event.Kind", event.Kind)
				switch event.Kind {
				case "zigbee_connectivity":
					// slient ignore connectivity update
				case "light":
					parts := strings.Split(event.Id_v1, "/")
					kind := parts[1]
					fmt.Println("kind", kind)

					if kind == "lights" {
						if len(parts[2]) == 0 {
							continue
						}

						var deviceId string
						hasSet := false

						for _, v := range s.devices {
							if v.id == parts[2] {
								deviceId = v.Uniqueid
							}
						}

						device := s.State[deviceId]
						evt := homeClient.NewEvent()
						if pVal, ok := event.Dimming["brightness"]; ok {

							fmt.Println("Bri", device.Bri, "= val", pVal)

							if device.Bri != int(pVal) {
								fmt.Println("add brightness", pVal)
								// TODO: need to convert %val to real number between 0-255

								device.Bri = int(pVal)
								val := int((pVal / 100) * 255)
								evt.AddDial("brightness", int(val))
								hasSet = true
							}
						}
						if val, ok := event.On["on"]; ok {
							fmt.Println("On", device.Bri, "= val", val)
							if device.On != val {
								fmt.Println("add state", val)
								device.On = val
								evt.AddSwitch("state", val)
								hasSet = true
							}
						}
						if val := event.Color; val != nil {
							fmt.Println("set Colour")
							// val.Xy.X
						}

						if hasSet {
							fmt.Println("client.SendEvent", deviceId, evt)
							s.State[deviceId] = device
							s.client.SendEvent(deviceId, evt)
						}
					}
				}
			}
		}
	}
}
