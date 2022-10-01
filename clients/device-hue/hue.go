package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	homeClient "device-hue/homeClient"
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
}

type hueConfig struct {
	Name     string
	Bridgeid string
}

func (s *settings) hueRegisterHub(username string, url string, client *homeClient.AhClient) {
	info := s.getHueInfo()

	hub := homeClient.NewHub(info.Config.Name, info.Config.Bridgeid)

	if len(s.devices) == 0 {
		s.devices = make(map[string]string)
	}

	//TODO: this dosent work device properties dont do anything
	for key, v := range info.Lights {
		dev := homeClient.NewDevice(v.Name, v.Uniqueid)

		dev.AddSwitch("state", v.Productname, v.State.On, "RW")
		dev.AddDial("brightness", v.Productname, v.State.Bri, 0, 254, "RW")
		hub.AddDevice(dev)

		s.devices[v.Uniqueid] = key
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
		return nil, err
	}

	return resp, nil
}
