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
	Token      string
	FifoFile   string
	ServerURL  string
	Deviceid   string
	DeviceName string
}

func main() {
	profile, err := os.UserConfigDir()
	if err != nil {
		log.Panic("unable to get users home folder", err)
	}
	configPath := path.Join(profile, "auto-home", "device.custom.json")

	jsonFile, err := os.Open(configPath)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var conf settings
	json.Unmarshal(byteValue, &conf)

	os.Remove(conf.FifoFile)

	client := homeClient.NewClient(conf.ServerURL, conf.Token)

	event, err := client.ListenEvents(conf.callback)
	if err != nil {
		log.Panic("unable to listen:", err)
	}

	dev := homeClient.NewDevice(conf.DeviceName, conf.Deviceid)

	dev.AddSwitch("state", "", "closed", "RW")

	client.RegisterDevice(&dev)

	err = syscall.Mkfifo(conf.FifoFile, 0666)
	if err != nil {
		log.Println("unable to create fifo", conf.FifoFile)
	}

	// to open pipe to read
	go func() {
		for {
			file, err := os.OpenFile(conf.FifoFile, os.O_RDONLY, os.ModeNamedPipe)
			if err != nil {
				log.Println("unable to open file", err)
			}

			scanner := bufio.NewScanner(file)
			// optionally, resize scanner's capacity for lines over 64K, see next example
			for scanner.Scan() {
				arrItem := strings.Split(scanner.Text(), "=")
				evt := homeClient.NewEvent()
				if arrItem[0] == "dial" {
					val, _ := strconv.Atoi(arrItem[2])
					evt.AddDial(arrItem[1], val)
				} else if (arrItem[0]) == "switch" {
					evt.AddSwitch(arrItem[1], arrItem[2])
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
			if msg == homeClient.EVENT_SHUTDOWN {
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
