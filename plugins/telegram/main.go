package main

import (
	"bufio"
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

// opena telegram connection to read any messages sent to the bot
func (t *Telegram) getMessage(accountName string) {
	account := t.conf[accountName]

	// Send the message
	// TODO: need to record and reuse the offset
	url := "https://api.telegram.org/" + account.BotID + ":" + account.Key + "/getUpdates" //?offset=972164103"

	client := http.Client{
		Timeout: 300 * time.Second,
	}

	for {
		// isConnected := false
		req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer([]byte{}))
		if err != nil {
			fmt.Println("eventstream read error:", err)
			break
		}
		res, err := client.Do(req)
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

		}
	}

	// Log
	log.Printf("listen connection closed")

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

	// TODO: this wont work how I want, as I cant match the users telegram id to their auto-home user id.
	//  until I work out how user messages are going to work I have disabled the below go func
	// go func() {
	// 	for name, _ := range conf {
	// 		fmt.Println(">> setting listener for account:", name)
	// 		cal.getMessage(name)
	// 	}
	// }()

	err = p.Done()
	if err != nil {
		fmt.Println(err)
	}

}
