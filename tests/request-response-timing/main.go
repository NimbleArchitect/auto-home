package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"
)

const registerDevices = true

const serverUrl = "https://localhost:4242/v1"

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

const regDevice = `{
	"method": "device",
	"data": {
		"id": "123-echo-321",
		"name": "echo device",
		"desc": "echos the recieved data",
		"help": "",
		"properties": [
			{
				"name": "switch",
				"description": "echo the switch state",
				"type": "switch",
				"value": "off",
				"mode": "rw"
			},{
				"name": "dial9",
				"description": "echo the switch state",
				"type": "dial",
				"min": 0,
				"max": 359,
				"value": 100,
				"mode": "rw"
			}
		]
	}
}
`

const deviceEvent = `{
	"Method": "event",
	"data": {
		"id": "123-echo-321",
		"properties": [
		{
			"name": "switch",
			"type": "switch",
			"value": "on"
		},{
			"name": "dial9",
			"type": "dial",
			"value": %VALUE%
		}
		],
		"timestamp": "20180415T142242"
	}
}
`

func main() {
	var sessionid string
	var actionid string
	var start time.Time

	pool, err := x509.SystemCertPool()
	if err != nil {
		log.Fatal(err)
	}

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

	hclient := &http.Client{
		Transport: roundTripper,
		Timeout:   time.Second * 600,
	}

	// login
	// connect and get a session id
	auth_data := []byte(fmt.Sprintf(`{"data":{"returnId": true, "user": "%s", "pass": "%s"}}`, "virtual.custom.light", "secretclientid"))
	req, err := http.NewRequest(http.MethodPost, serverUrl+"/connect", bytes.NewBuffer([]byte(auth_data)))
	if err != nil {
		log.Println(err)
		return
	}

	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := hclient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	type AuthResult struct {
		Result jsonStatus
		Data   map[string]string
	}

	var result AuthResult
	err = json.NewDecoder(resp.Body).Decode(&result)

	if result.Result.Status != "ok" {
		fmt.Println("unable to connect")
		return
	}

	data := result.Data
	if val, ok := data["session"]; ok {
		sessionid = val
	} else {
		fmt.Println("invalid session")
	}
	if val, ok := data["actionid"]; ok {
		actionid = val
	} else {
		fmt.Println("invalid session")
	}

	if registerDevices {
		// registration
		// add devices to system
		json_data := []byte(regDevice)

		// log.Println("register device")

		req, err = http.NewRequest(http.MethodPost, serverUrl+"/register", bytes.NewBuffer(json_data))
		if err != nil {
			log.Println(err)
			return
		}

		// set the request header Content-Type for json
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		req.Header.Set("session", sessionid)
		r, err := hclient.Do(req)
		if err != nil {
			log.Println(err)
			return
		}

		var tmp Result
		// log.Println("decode json")
		err = json.NewDecoder(r.Body).Decode(&tmp)
		if err != nil {
			log.Println(err)
		}

		// fmt.Println(tmp)
		if tmp.Result.Status != "ok" {
			log.Println("ERROR: empty response from register call")
			os.Exit(1)
		}
	}

	// fmt.Println(">>", tmp)

	ready := make(chan bool, 1)

	var timeing []time.Duration

	go func() {
		log.Println("starting scanner")
		for {
			req, err = http.NewRequest(http.MethodGet, serverUrl+"/actions/"+actionid, bytes.NewBuffer([]byte("")))
			if err != nil {
				log.Println(err)
				return
			}

			// set the request header Content-Type for json
			req.Header.Set("Content-Type", "application/json; charset=utf-8")
			if len(sessionid) > 0 {
				req.Header.Set("session", sessionid)
			}

			out, err := hclient.Do(req)
			if err != nil {
				log.Println(err)
				return
			}

			scanner := bufio.NewScanner(out.Body)
			if scanner.Err() != nil {
				log.Println("unable to start scanner")
				break
			}
			for scanner.Scan() {
				// on a read error this loop breaks out

				ln := scanner.Text()
				if len(ln) == 0 {
					continue
				}
				// log.Println("scanner recieved", ln)
				if ln == "{\"Method\": \"shutdown\"}" {
					break
				}
				elapsed := time.Since(start)
				fmt.Printf("time taken %s\n", elapsed)

				var tmp actionResult

				json.Unmarshal([]byte(ln), &tmp)
				if tmp.Method == "action" {
					for _, v := range tmp.Data.Properties {
						if v["name"] == "dial9" {
							f, _ := v["value"].(float64)
							s := strconv.FormatFloat(f, 'f', 0, 32)
							val, _ := strconv.Atoi(s)
							if val == 9 {
								timeing = append(timeing, elapsed)
							} else {
								log.Fatal("something went wrong:", ln)
							}
						}
					}
				}

				// done <- true
				ready <- true
			}
			// log.Println("scanner dropped:", scanner.Err())

			// log.Println("re-starting scanner")
			time.Sleep(2 * time.Second)
			// os.Exit(1)
		}
	}()

	time.Sleep(1 * time.Second)

	ready <- true

	for i := 1; i <= 10; i = i + 1 {

		start = time.Now()

		eventurl := serverUrl + "/event/" + actionid

		json_data := []byte(strings.Replace(deviceEvent, "%VALUE%", strconv.Itoa(i), 1))

		// json_data = getDeviceJson(jsonevent, i)

		// log.Println("Post event:", eventurl)

		req, err = http.NewRequest(http.MethodPost, eventurl, bytes.NewBuffer(json_data))
		if err != nil {
			log.Println(err)
			return
		}

		// set the request header Content-Type for json
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		req.Header.Set("session", sessionid)
		r, err := hclient.Do(req)

		if err != nil {
			log.Fatal(err)
		}

		msg, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}

		_ = msg

		<-ready
	}

	var max time.Duration
	for _, v := range timeing {
		if v > max {
			max = v
		}
	}

	fmt.Println("longest", max)

	if max <= 10000*time.Nanosecond {
		os.Exit(1)
	}

	if max > 5*time.Millisecond {
		os.Exit(1)
	}

	os.Exit(0)

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
