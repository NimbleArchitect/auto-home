package homeClient

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"
)

const (
	FLAG_BUTTON = iota
	FLAG_DIAL
	FLAG_SWITCH
	FLAG_TEXT
)

const (
	EVENT_RESTART = iota
	EVENT_SHUTDOWN
	EVENT_RELOAD
)

type Result struct {
	Result jsonStatus
	Data   map[string]string
}

type jsonStatus struct {
	Status string
	Msg    string
}

type actionResult struct {
	Method string
	Data   actionData
}

type actionData struct {
	ID         string
	Properties []map[string]interface{}
}

type AhClient struct {
	http      *http.Client
	uuid      string
	done      chan bool
	address   string
	sessionId string
}

func NewClient(address string, token string) AhClient {
	// var start time.Time
	// keyLogFile := "./key.log"

	// var keyLog io.Writer
	// if len(keyLogFile) > 0 {
	// 	f, err := os.Create(keyLogFile)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	defer f.Close()
	// 	keyLog = f
	// }

	pool, err := x509.SystemCertPool()
	if err != nil {
		log.Fatal(err)
	}
	// AddRootCA(pool)

	qconf := quic.Config{
		KeepAlivePeriod: 60 * time.Second,
		MaxIdleTimeout:  600 * time.Second,
	}

	roundTripper := &http3.RoundTripper{
		TLSClientConfig: &tls.Config{
			RootCAs:            pool,
			InsecureSkipVerify: true,
			// KeyLogWriter:       keyLog,
		},
		QuicConfig: &qconf,
	}

	// hclient := &http.Client{
	// 	Transport: roundTripper,
	// 	Timeout:   time.Second * 600,
	// }
	out := AhClient{
		address: address,
		http: &http.Client{
			Transport: roundTripper,
			Timeout:   time.Second * 600,
		},
		sessionId: "",
	}

	out.Connect(token)
	return out
}

func (c *AhClient) Connect(token string) {
	data := `{"token":"` + token + `"}`

	r, err := c.makeRequest(c.address+"/connect", http.MethodPost, data)

	if err != nil {
		log.Println(err)
		return
	}

	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(">>", r.StatusCode)
	c.sessionId = r.Header.Get("session")
	fmt.Println(">> session:", c.sessionId)
}

func (c *AhClient) RegisterDevice(device *Device) {

	jsonData := device.getJson()
	jsonOut := fmt.Sprintf("{\"method\":\"device\",\"data\":%s}", jsonData)

	log.Println("register device")
	fmt.Println(">> jsonOut", jsonOut)
	r, err := c.makeRequest(c.address+"/register", http.MethodPost, jsonOut)

	if err != nil {
		log.Fatal(err)
	}

	var tmp Result
	log.Println("decode json")
	err = json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		log.Println(err)
	}

	if tmp.Result.Status != "ok" {
		log.Println("ERROR: empty response from register call")
		os.Exit(1)
	}

}

func (c *AhClient) RegisterHub(hub *Hub) {

	var deviceListJson string
	for _, v := range hub.data.Devices {
		deviceListJson += v.getJson() + ","
	}

	deviceListJson = deviceListJson[0 : len(deviceListJson)-1]

	jsonData := fmt.Sprintf("\"id\":\"%s\",\"name\":\"%s\",\"description\":\"%s\",\"devices\":[%s]", hub.data.Id, hub.data.Name, hub.data.Description, deviceListJson)
	jsonOut := fmt.Sprintf("{\"method\":\"hub\",\"data\":{%s}}", jsonData)

	// println(jsonOut)

	log.Println("register device")
	r, err := c.makeRequest(c.address+"/register", http.MethodPost, jsonOut)

	if err != nil {
		log.Fatal(err)
	}

	if r.StatusCode == 500 {
		fmt.Println(">> httpStatus:", r.Status)
		os.Exit(1)
	}

	var tmp Result
	log.Println("decode json")
	err = json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		log.Println("decodeError:", err)
	}

	if tmp.Result.Status != "ok" {
		log.Println("ERROR: empty response from register call")
		os.Exit(1)
	}

}

func (c *AhClient) Close() {
	c.done <- true
}

// ListenEvents starts the listener and calls the provided callback function on recieved events,
//
//	the returned channel allows the caller to respond to EVENT flags
func (c *AhClient) ListenEvents(callback func(string, map[string]interface{})) (chan int, error) {
	if c.done != nil {
		return nil, errors.New("already listening")
	}
	c.done = make(chan bool, 1)

	eventAction := make(chan int)

	ready := make(chan bool)

	// first call the listener to get our id
	json_data := `{"method": "listen"}`
	r, err := c.makeRequest(c.address+"/register", http.MethodPost, json_data)
	if err != nil {
		return nil, err
	}

	var tmp Result
	log.Println("decode json")
	err = json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		return nil, err
	}

	if tmp.Result.Status != "ok" {
		return nil, errors.New("ERROR: empty response from listen call")

	}

	c.uuid = tmp.Data["id"]

	// now we can start listening
	go c.startListener(&ready, eventAction, callback)
	<-ready
	return eventAction, nil
}

func (c *AhClient) startListener(ready *chan bool, eventAction chan int, callback func(string, map[string]interface{})) {
	go func() {
		log.Println("starting scanner")
		for {
			fmt.Println("connecting to", c.address+"/actions/"+c.uuid)
			out, err := c.makeRequest(c.address+"/actions/"+c.uuid, http.MethodGet, "")

			if err != nil {
				fmt.Println("unable to connect", err)
			}

			scanner := bufio.NewScanner(out.Body)
			if scanner.Err() != nil {
				log.Println("unable to start scanner")
				break
			}
			if ready != nil {
				*ready <- true
			}

			for scanner.Scan() {
				// on a read error this loop breaks out

				ln := scanner.Text()
				if len(ln) == 0 {
					//recieved empty line from the server
					// treat it as a connection test and ignore
					continue
				}

				log.Println("scanner recieved", ln)
				var tmp actionResult
				log.Println("decode json")

				json.Unmarshal([]byte(ln), &tmp)
				switch tmp.Method {
				case "action":
					props := make(map[string]interface{})
					for _, v := range tmp.Data.Properties {
						name := v["name"].(string)
						val := v["value"]
						props[name] = val
					}
					callback(tmp.Data.ID, props)

				case "shutdown":
					eventAction <- EVENT_SHUTDOWN

				case "restart":
					eventAction <- EVENT_RESTART

				case "reload":
					eventAction <- EVENT_RESTART

				default:
					log.Println("recieved unknown method:", tmp.Method)
				}

				if err != nil {
					log.Println(err)
				}

				// done <- true
			}
			log.Println("scanner dropped:", scanner.Err())

			log.Println("re-starting scanner")
			// time.Sleep(1 * time.Second)
		}
	}()

	<-c.done
	fmt.Println(">> listen finished")
}

func (c *AhClient) SendEvent(deviceid string, evt event) error {
	if len(c.sessionId) == 0 {
		return errors.New("invalid session, you must connect first")
	}

	eventurl := c.address + "/event/" + c.uuid

	var propJson string
	for _, v := range evt.props {
		propJson += v.json + ","
	}

	propJson = propJson[0 : len(propJson)-2]

	fmt.Println(">> propJson", propJson)

	jsonData := fmt.Sprintf("\"id\":\"%s\",\"properties\":[%s]", deviceid, propJson)
	jsonOut := fmt.Sprintf("{\"Method\":\"event\",\"data\":{%s}}", jsonData)

	log.Println("Post event:", eventurl)
	r, err := c.makeRequest(eventurl, http.MethodPost, jsonOut)

	if err != nil {
		return err
	}

	msg, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	log.Println(string(r.Status), "++", string(msg))
	return nil
}

func (c *AhClient) makeRequest(url string, method string, data string) (*http.Response, error) {

	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return nil, err
	}

	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	if len(c.sessionId) > 0 {
		req.Header.Set("session", c.sessionId)
	}
	r, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func AddRootCA(certPool *x509.CertPool) {
	caCertPath := "./ca.pem"
	caCertRaw, err := os.ReadFile(caCertPath)
	if err != nil {
		panic(err)
	}
	if ok := certPool.AppendCertsFromPEM(caCertRaw); !ok {
		panic("Could not add root ceritificate to pool.")
	}
}
