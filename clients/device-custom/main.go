package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"syscall"

	homeClient "device-custom/homeClient"
)

type settings struct {
	Token         string
	FifoFile      string
	ServerURL     string
	Deviceid      string
	DeviceName    string
	Clientid      string
	PropertyName  string
	PropertyType  string
	PropertyState string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("you must provide the configuration filename as an argument")
		os.Exit(1)
	}
	configFile := os.Args[1]

	profile, err := os.UserConfigDir()
	if err != nil {
		log.Panic("unable to get users home folder", err)
	}
	configPath := path.Join(profile, "auto-home", configFile)

	jsonFile, err := os.Open(configPath)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var conf settings
	json.Unmarshal(byteValue, &conf)

	os.Remove(conf.FifoFile)

	client := homeClient.NewClient(conf.ServerURL, conf.Clientid, conf.Token)

	event, err := client.ListenEvents(conf.callback)
	if err != nil {
		log.Panic("unable to listen:", err)
	}

	dev := homeClient.NewDevice(conf.DeviceName, conf.Deviceid)

	switch conf.PropertyType {
	case "switch":
		dev.AddSwitch(conf.PropertyName, "", conf.PropertyState, "RW")
	case "button":
		if conf.PropertyState == "true" {
			dev.AddButton(conf.PropertyName, "", true, true, "RW")
		} else {
			dev.AddButton(conf.PropertyName, "", false, false, "RW")
		}
	}

	client.RegisterDevice(&dev)

	err = syscall.Mkfifo(conf.FifoFile, 0666)
	if err != nil {
		log.Println("unable to create fifo", conf.FifoFile)
	}

	go func() {
		for {
			// to open pipe to read
			file, err := os.OpenFile(conf.FifoFile, os.O_RDONLY, os.ModeNamedPipe)
			if err != nil {
				log.Println("unable to open file", err)
			}

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				arrItem := strings.Split(scanner.Text(), "=")
				evt := homeClient.NewEvent()
				if arrItem[0] == "dial" {
					val, _ := strconv.Atoi(arrItem[2])
					evt.AddDial(arrItem[1], val)
				} else if (arrItem[0]) == "switch" {
					evt.AddSwitch(arrItem[1], arrItem[2])
				} else if (arrItem[0]) == "button" {
					if arrItem[2] == "true" {
						evt.AddButton(arrItem[1], true)
					} else {
						evt.AddButton(arrItem[1], false)
					}
				}
				client.SendEvent(conf.Deviceid, evt)
				fmt.Println(">> event sent!")
				break
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}

			if file != nil {
				file.Close()
			}
		}
	}()

	finished := false
	for {
		select {
		case msg := <-event:
			switch msg {
			case homeClient.EVENT_SHUTDOWN:
				finished = true
			}
		}
		if finished {
			break
		}
	}

	fmt.Println("!>> got here")

}

func (s *settings) callback(deviceid string, args map[string]interface{}) {

}
