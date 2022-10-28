package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"log"
	"net/http"
	"os"
	"path"
)

// const SockAddr = "/tmp/rpc.sock"

type Client int

type Telegram struct {
	conf settings
}

type settings struct {
	BotID    string
	Key      string
	ChatId   string
	SockAddr string
}

type Result struct {
	Ok   bool
	Data map[string]interface{}
}

func (t *Telegram) SendMessage(raw []byte) {
	// var m event
	var val map[int]interface{}

	err := json.Unmarshal(raw, &val)
	if err != nil {
		fmt.Println("e>>", err)
	}

	msg := val[0].(string)

	fmt.Println("sending telegram>>", msg)

	// Global variables
	var response *http.Response

	// Send the message
	url := "https://api.telegram.org/" + t.conf.BotID + ":" + t.conf.Key + "/sendMessage"
	body, _ := json.Marshal(map[string]string{
		"chat_id": t.conf.ChatId,
		"text":    msg,
	})

	client := http.Client{
		Timeout: 1 * time.Second,
	}

	response, err = client.Post(
		url,
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return
	}

	// Close the request at the end
	defer response.Body.Close()

	// Body
	body, err = io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("err<<", err)
		return
	}

	// Log
	log.Printf("Message '%s' was sent", msg)
	log.Printf("Response JSON: %s", string(body))

}

type event struct {
	Label    string
	Date     time.Time
	Location string
	Notes    string
}

func main() {

	profile, err := os.UserConfigDir()
	if err != nil {
		log.Panic("unable to get users home folder", err)
	}
	configPath := path.Join(profile, "auto-home", "plugin.telegram.json")

	jsonFile, err := os.Open(configPath)
	if err != nil {
		log.Println("unable to open plugin.telegram.json", err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var conf settings
	json.Unmarshal(byteValue, &conf)

	p := Connect(SockAddr)

	cal := new(Telegram)
	cal.conf = conf
	p.Register("telegram", cal)

	// go func() {
	ev := event{
		Label:    "car Mot",
		Date:     time.Now(),
		Location: "home",
		Notes:    "",
	}

	time.Sleep(10 * time.Second)
	p.Call("onEvent", ev)
	// }()

	err = p.Done()
	if err != nil {
		fmt.Println(err)
	}

}
