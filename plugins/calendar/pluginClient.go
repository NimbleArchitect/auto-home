package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"reflect"
	"sync"
	"time"
	"unicode"
)

type Generic struct {
	Method string
	Id     int
	Data   *json.RawMessage
}

type create struct {
	Name   string
	Fields map[string]interface{}
}

type request struct {
	Method string
	Id     int
	Data   interface{}
}

type trigger struct {
	Name   string
	Call   string
	Fields interface{}
}

type action struct {
	Name   string
	Call   string
	Fields *json.RawMessage
}

type Args struct {
	Kind  string
	Value interface{}
}

type plugin struct {
	c          connector
	lock       sync.Mutex
	chError    chan error
	functions  map[string]reflect.Value
	name       string
	isFinished chan error
}

func (p *plugin) Call(callName string, arg interface{}) {
	nextId := p.c.WaitAdd()

	generic := request{
		Method: "trigger",
		Id:     nextId,
		Data: trigger{
			Name:   p.name,
			Call:   callName,
			Fields: arg,
		},
	}

	data, err := json.Marshal(generic)
	if err != nil {
		log.Println("json error", err)
		// resp := p.MakeError(generic.Id, err)
		// p.c.c.SetWriteDeadline(time.Now().Add(5 * time.Second))
		// p.c.WriteB(resp)
	}

	p.c.WriteB(data)
	p.c.WaitOn(nextId)

}

func (p *plugin) Register(name string, obj interface{}) {
	//register an object and its methods
	methods := make(map[string]interface{})
	p.name = name
	p.c.name = &p.name

	nextId := p.c.WaitAdd()

	t := reflect.TypeOf(obj)
	for i := 0; i < t.NumMethod(); i++ {
		rawName := t.Method(i).Name
		name := []rune(rawName)
		name[0] = unicode.ToLower(name[0])
		function := reflect.ValueOf(obj).MethodByName(rawName)
		p.functions[string(name)] = function
		methods[string(name)] = nil
	}

	generic := request{
		Method: "create",
		Id:     nextId,
		Data: create{
			Name:   p.name,
			Fields: methods,
		},
	}

	data, err := json.Marshal(generic)
	if err != nil {
		log.Println("json error", err)
		// resp := p.MakeError(generic.Id, err)
		// p.c.c.SetWriteDeadline(time.Now().Add(5 * time.Second))
		// p.c.WriteB(resp)
	}

	p.c.WriteB(data)
	p.c.WaitOn(nextId)

	// if _, ok := p.c.wait[p.c.NextId]; !ok {
	// 	p.c.wait[p.c.NextId] = make(chan bool, 1)
	// 	p.c.WriteB(data)
	// 	<-p.c.wait[p.c.NextId]
	// }

}

// Done wait untile the server tells the plugin to exit
func (p *plugin) Done() error {
	err := <-p.isFinished
	return err
}

func Connect() *plugin {
	addr, ok := os.LookupEnv("autohome_sockaddr")
	if !ok {
		log.Println("unable to get host address, are you running from within auto-home")
		os.Exit(1)
	}

	conn, err := net.Dial("unix", addr)
	if err != nil {
		panic(err)
	}
	// defer conn.Close()

	out := plugin{
		c: connector{
			c:            conn,
			lock:         sync.Mutex{},
			nextId:       0,
			responseWait: make(map[int]*chan result),
		},
		lock:       sync.Mutex{},
		chError:    make(chan error),
		functions:  make(map[string]reflect.Value),
		isFinished: make(chan error, 2),
	}

	go out.handle()
	return &out
}

func (p *plugin) handle() {
	var buf []byte

	defer p.c.c.Close()

	go func() {

		for {
			tmp := make([]byte, 256)
			n, err := p.c.c.Read(tmp)

			if err != nil {
				p.chError <- err
				break
			} else {
				buf = append(buf, tmp[:n]...)
				l := len(buf)

				if buf[l-2] == 10 && buf[l-1] == 10 {
					if len(buf) == 2 {
						buf = []byte{}
						continue
					}
					last := 0
					for i := 1; i < len(buf); i++ {
						if buf[i-1] == 10 && buf[i] == 10 {
							fmt.Println(p.name, "=<< recieved", string(buf[last:i-2]))
							go p.decode(buf[last:i])
							last = i
						}

					}
					buf = []byte{}
				}
			}
		}
	}()

	for {
		select {
		case err := <-p.chError:
			if io.EOF == err {
				if len(buf) > 0 {
					p.decode(buf)
				}
				p.isFinished <- err
				return
			}
		case <-time.After(1 * time.Minute):
			p.c.WriteB([]byte(nil))

		}
	}

}

func (p *plugin) decode(buf []byte) {
	var err error
	var generic Generic

	err = json.Unmarshal(buf, &generic)
	if err != nil {
		log.Println("decode error", err)
		resp := p.MakeError(generic.Id, err)
		p.c.WriteB(resp)
	} else {
		//process message
		p.processMessage(generic)
	}

}

func (p *plugin) processMessage(obj Generic) error {
	raw, _ := obj.Data.MarshalJSON()
	switch obj.Method {

	case "result":
		var m result
		p.c.WaitDone(obj.Id, nil, nil)
		json.Unmarshal(raw, &m)

	case "trigger":
		// call the methods that where regestered from the object
		var m action
		json.Unmarshal(raw, &m)

		raw, err := m.Fields.MarshalJSON()
		if err != nil {
			return err
		}

		fmt.Println("raw:", string(raw))

		var callArgs []reflect.Value
		var response []reflect.Value

		function := p.functions[m.Call]
		if len(raw) <= 2 {
			response = function.Call(nil)
		} else {
			v := reflect.ValueOf(raw)
			callArgs = []reflect.Value{v}

			// TODO: callArgs has to be a string, this means I need a way to convert between the json string and the actual struct
			response = function.Call(callArgs)
		}

		var retValues interface{}
		if len(response) > 0 {
			// process answer from function.call
			// TODO: response is an array so need to process each value

			if response[0].IsValid() {
				retValues = response[0].Interface()
			}

		}

		out := p.MakeResponse(obj.Id, retValues)

		p.c.WriteB(out)
		p.c.WaitDone(obj.Id, nil, nil)
	}

	return nil
}

type Response struct {
	Method string
	Id     int
	Data   interface{}
}

type result struct {
	Ok      bool
	Message string                  `json:",omitempty"`
	Data    *map[string]interface{} `json:",omitempty"`
}

func (p *plugin) MakeError(id int, err error) []byte {
	var msg string

	ok := false
	if err == nil {
		ok = true
	} else {
		msg = err.Error()
	}

	resp := Response{
		Method: "result",
		Id:     id,
		Data: result{
			Ok:      ok,
			Message: msg,
		},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		log.Fatalln("json error", err)
		// if we cant encode, we can't comunicate with the server so
		//  we bail with an error
		p.isFinished <- err
	}

	return data
}

func (p *plugin) MakeResponse(id int, singleArg interface{}) []byte {
	arg := make(map[string]interface{})
	arg["0"] = singleArg

	resp := Response{
		Method: "result",
		Id:     id,
		Data: result{
			Ok:   true,
			Data: &arg,
		},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		log.Println("json error", err)
	}

	return data
}

type connector struct {
	c            net.Conn
	lock         sync.Mutex
	nextId       int
	responseWait map[int]*chan result
	name         *string
}

func (c *connector) WaitAdd() int {
	i := c.nextId
	c.nextId++

	wait := make(chan result, 1)
	c.lock.Lock()
	c.responseWait[i] = &wait
	c.lock.Unlock()

	return i
}

// WaitDone used to signal that we should not wait any more
func (c *connector) WaitDone(i int, msg *string, data *map[string]interface{}) {
	var msgData string

	c.lock.Lock()
	wait, ok := c.responseWait[i]
	if !ok {
		channel := make(chan result, 1)
		wait = &channel
		c.responseWait[i] = wait
	}
	c.lock.Unlock()

	if msg != nil {
		msgData = *msg
	}
	*wait <- result{Ok: true, Message: msgData, Data: data}
}

func (c *connector) WaitOn(i int) (string, map[string]interface{}, bool) {
	c.lock.Lock()
	wait := *c.responseWait[i]
	c.lock.Unlock()

	out := <-wait

	var args map[string]interface{}

	return out.Message, args, out.Ok
}

func (c *connector) WriteB(b []byte) {
	//fmt.Println(*c.name, "=>> sending", string(b))
	c.lock.Lock()
	_, err := c.c.Write(b)
	if err != nil {
		log.Println("writeB error:", err)
	}
	c.c.Write([]byte("\n\n"))
	c.lock.Unlock()
}
