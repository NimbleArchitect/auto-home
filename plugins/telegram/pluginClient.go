package main

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"reflect"
	"sync"
	"time"
	"unicode"
)

const SockAddr = "/tmp/rpc.sock"

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
	p.c.NextId++

	generic := request{
		Method: "trigger",
		Id:     p.c.NextId,
		Data: trigger{
			Name:   p.name,
			Call:   callName,
			Fields: arg,
		},
	}

	data, err := json.Marshal(generic)
	if err != nil {
		// log.Println("json error", err)
		resp := p.MakeError(generic.Id, err)
		p.c.c.SetWriteDeadline(time.Now().Add(5 * time.Second))
		p.c.WriteB(resp)
	}

	if _, ok := p.c.wait[p.c.NextId]; !ok {
		p.c.wait[p.c.NextId] = make(chan bool, 1)
		p.c.WriteB(data)
		<-p.c.wait[p.c.NextId]
	}

}

func (p *plugin) Register(name string, obj interface{}) {
	//register an object and its methods
	methods := make(map[string]interface{})
	p.name = name

	p.c.NextId++

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
		Id:     p.c.NextId,
		Data: create{
			Name:   p.name,
			Fields: methods,
		},
	}

	data, err := json.Marshal(generic)
	if err != nil {
		// log.Println("json error", err)
		resp := p.MakeError(generic.Id, err)
		p.c.c.SetWriteDeadline(time.Now().Add(5 * time.Second))
		p.c.WriteB(resp)
	}

	if _, ok := p.c.wait[p.c.NextId]; !ok {
		p.c.wait[p.c.NextId] = make(chan bool, 1)
		p.c.WriteB(data)
		<-p.c.wait[p.c.NextId]
	}

}

func (p *plugin) Done() error {
	err := <-p.isFinished
	return err
}

func Connect(addr string) *plugin {
	conn, err := net.Dial("unix", addr)
	if err != nil {
		panic(err)
	}
	// defer conn.Close()

	out := plugin{
		c: connector{
			c:      conn,
			lock:   sync.Mutex{},
			NextId: -1,
			wait:   make(map[int]chan bool),
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
							// fmt.Println("--<< recieved", string(buf[last:i-2]))
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
			p.c.c.SetWriteDeadline(time.Now().Add(2 * time.Second))
			p.c.WriteB([]byte(nil))

		}
	}

}

func (p *plugin) decode(buf []byte) {
	var err error
	var generic Generic

	err = json.Unmarshal(buf, &generic)
	if err != nil {
		// log.Println("decode error", err)
		resp := p.MakeError(generic.Id, err)
		p.c.c.SetWriteDeadline(time.Now().Add(5 * time.Second))
		p.c.WriteB(resp)

	} else {
		// fmt.Println(">>", generic.Method)
		//process message
		p.processMessage(generic)
	}

}

func (p *plugin) processMessage(obj Generic) error {
	raw, _ := obj.Data.MarshalJSON()
	switch obj.Method {

	case "result":
		var m result
		json.Unmarshal(raw, &m)
		p.c.lock.Lock()
		p.c.wait[obj.Id] <- true
		p.c.lock.Unlock()

	case "trigger":
		// call the methods that where regestered from the object
		var m action
		json.Unmarshal(raw, &m)

		raw, err := m.Fields.MarshalJSON()
		if err != nil {
			return err
		}
		v := reflect.ValueOf(raw)
		function := p.functions[m.Call]
		function.Call([]reflect.Value{v})

		resp := p.MakeError(obj.Id, err)
		p.c.c.SetWriteDeadline(time.Now().Add(5 * time.Second))

		p.c.WriteB(resp)
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
	Call    string
	Message string
}

func (p *plugin) MakeError(id int, err error) []byte {
	var ok bool
	var msg string

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

type connector struct {
	c      net.Conn
	lock   sync.Mutex
	NextId int
	wait   map[int]chan bool
}

func (c *connector) WriteB(b []byte) {
	// fmt.Println("->> sending", string(b))
	c.lock.Lock()
	c.c.Write(b)
	c.c.Write([]byte("\n\n"))
	c.lock.Unlock()
}
