package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"path"
	"time"
)

const SockAddr = "/tmp/rpc.sock"

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

// This procedure is invoked by rpc and calls rpcexample.Multiply which stores product of args.A and args.B in result pointer
func (t *Client) RoleCall(args map[string]interface{}, result *Result) error {

	result.Data = make(map[string]interface{})

	result.Data["name"] = "Telegram"
	result.Data["SendMessage"] = ""
	result.Ok = true

	return nil
}

func (t *Telegram) SendMessage(args map[string]interface{}, result *Result) error {

	result.Data = make(map[string]interface{})

	// fmt.Println("SendMessage called")
	msg := args["message"].(string)

	// fmt.Println(msg)

	// Global variables
	var err error
	var response *http.Response

	// time.Sleep(5 * time.Second)
	// Send the message
	url := "https://api.telegram.org/" + t.conf.BotID + ":" + t.conf.Key + "/sendMessage"
	body, _ := json.Marshal(map[string]string{
		"chat_id": t.conf.ChatId,
		"text":    msg,
	})
	response, err = http.Post(
		url,
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		// fmt.Println(err)
		result.Data["error"] = err.Error()
		return nil
	}

	// Close the request at the end
	defer response.Body.Close()

	// Body
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("err<<", err)
		result.Data["error"] = err.Error()
		return nil
	}

	// Log
	log.Printf("Message '%s' was sent", msg)
	log.Printf("Response JSON: %s", string(body))

	result.Data["msg"] = string(body)
	result.Ok = true

	return nil
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

	reg := new(Client)
	arith := new(Telegram)
	arith.conf = conf

	conn, err := net.Dial("unix", SockAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	svr := rpc.NewServer()
	svr.Register(arith)
	svr.Register(reg)
	fmt.Println("telegram connected")
	svr.ServeConn(conn)

	time.Sleep(5 * time.Second)
}
