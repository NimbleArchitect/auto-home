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
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"
)

const eventCount = 402
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
				"name": "dialdelay",
				"description": "delayed dial",
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
			"name": "dialdelay",
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
	req, err := http.NewRequest(http.MethodPost, serverUrl+"/connect", bytes.NewBuffer(auth_data))
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

	expectedValue := 0
	sentValues := 0
	matchCount := 0

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
			req.Header.Set("session", sessionid)

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
				// log.Println("scanner recieved", ln)
				if ln == "{\"Method\": \"shutdown\"}" {
					break
				}

				var tmp actionResult

				json.Unmarshal([]byte(ln), &tmp)
				if tmp.Method == "action" {
					for _, v := range tmp.Data.Properties {
						if v["name"] == "dialout" {
							f, _ := v["value"].(float64)
							s := strconv.FormatFloat(f, 'f', 0, 32)
							val, _ := strconv.Atoi(s)

							// val is the incoming value
							if val == expectedValue {
								matchCount++
							} else {
								fmt.Println("incorrect value recieved")
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

	wg := sync.WaitGroup{}
	wg.Add(1)

	s2 := rand.NewSource(42)
	r2 := rand.New(s2)
	// fmt.Print(r2.Intn(100), ",")
	start := time.Now()

	for i := 1; i <= 402; i = i + 1 {

		wg.Add(1)

		eventurl := serverUrl + "/event/" + actionid

		ri := r2.Intn(300)
		expectedValue = ri
		sentValues++

		// go func() {
		json_data := []byte(strings.Replace(deviceEvent, "%VALUE%", strconv.Itoa(ri), 1))

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
		wg.Done()
		_ = msg
		// }()
	}

	wg.Done()
	wg.Wait()

	elapsed := time.Since(start)

	////////////////////////////////////////////////////////////
	// if sentValues != matchCount {
	// 	log.Panicln("sent count does not match recieved count")
	// 	os.Exit(1)
	// } else {
	// 	fmt.Println("sent count matches recieve count")
	// }

	fmt.Printf("sent %d events in %f seconds\n", eventCount, elapsed.Seconds())
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
