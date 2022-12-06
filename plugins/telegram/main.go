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

type Client int

type Telegram struct {
	conf map[string]settings
}

type settings struct {
	Name   string
	BotID  string
	Key    string
	ChatId string
}

func (t *Telegram) SendMessage(raw []byte) {
	// var m event
	var val map[int]interface{}
	var accountName string

	err := json.Unmarshal(raw, &val)
	if err != nil {
		fmt.Println("e>>", err)
	}

	args := val[0].(map[string]interface{})
	msg := args["message"].(string)
	name, ok := args["account"]
	if ok {
		accountName = name.(string)
	}
	if len(accountName) <= 0 {
		accountName = "default"
	}

	account := t.conf[accountName]
	// fmt.Println("sending telegram>>", msg)

	// Global variables
	var response *http.Response

	// Send the message
	url := "https://api.telegram.org/" + account.BotID + ":" + account.Key + "/sendMessage"
	body, _ := json.Marshal(map[string]string{
		"chat_id": account.ChatId,
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

	var onDiskSettings []settings
	json.Unmarshal(byteValue, &onDiskSettings)

	conf := make(map[string]settings)
	for _, v := range onDiskSettings {
		conf[v.Name] = v
	}

	p := Connect()

	cal := new(Telegram)
	cal.conf = conf
	p.Register("telegram", cal)

	// go func() {
	// 	ev := event{
	// 		Label:    "car Mot",
	// 		Date:     time.Now(),
	// 		Location: "home",
	// 		Notes:    "",
	// 	}

	// 	time.Sleep(4 * time.Second)

	// 	// TODO: callback dosent work, looks like im not recieving the result or its not being sent needs more investigation
	// 	p.Call("onEvent", ev)
	// }()

	err = p.Done()
	if err != nil {
		fmt.Println(err)
	}

}
