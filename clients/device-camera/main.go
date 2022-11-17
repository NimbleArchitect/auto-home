package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"

	homeClient "device-camera/homeClient"
)

type settings struct {
	Username   string
	HubAddress string
	FifoFile   string
	ServerURL  string
	// devices    map[string]string
	// http       *http.Client
}

func echoServer(c net.Conn) {
	log.Printf("Client connected [%s]", c.RemoteAddr().Network())
	io.Copy(c, c)
	c.Close()
}

func main() {
	const token = "randomcameradeviceuuid"

	profile, err := os.UserConfigDir()
	if err != nil {
		log.Panic("unable to get users home folder", err)
	}
	configPath := path.Join(profile, "auto-home", "device.camera.json")

	jsonFile, err := os.Open(configPath)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var conf settings
	json.Unmarshal(byteValue, &conf)

	client := homeClient.NewClient(conf.ServerURL, "virtual.custom.camera", token)

	event, err := client.ListenEvents(conf.callback)
	if err != nil {
		log.Panic("unable to listen", err)
	}

	dev := homeClient.NewDevice("Garden Camaer", "garden-cam")

	dev.AddSwitch("record", "", false, "RW")
	// dev.AddDial("event", "", false, 0, 254, "RW")
	client.RegisterDevice(&dev)

	err = syscall.Mkfifo(conf.FifoFile, 0666)
	if err != nil {
		log.Println("unable to create fifo", conf.FifoFile)
	}
	// to open pipe to write
	// file, err1 := os.OpenFile("tmpPipe", os.O_RDWR, os.ModeNamedPipe)

	// to open pipe to read
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
				// val, _ := strconv.ParseBool()
				evt.AddSwitch(arrItem[1], arrItem[2])
			}
			client.SendEvent("garden-cam", evt)
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		if file != nil {
			file.Close()
		}
	}

	finished := false
	for {
		select {
		case msg := <-event:
			switch msg {
			case homeClient.EVENT_SHUTDOWN:
				finished = true
			}
		case <-time.After(10 * time.Second):
			log.Println(">> pull state <<")
		}
		if finished {
			break
		}
	}

	fmt.Println("!>> got here")
	// evt := homeClient.NewEvent()

	// evt.AddDial("hue", 50)
	// evt.AddSwitch("state", "on")

	// client.SendEvent("123-tv-light-321", evt)

	// time.Sleep(600 * time.Second)
}

func (s *settings) callback(deviceid string, args map[string]interface{}) {
	// var j string

	// fmt.Println(">>", s.HubAddress)

}
