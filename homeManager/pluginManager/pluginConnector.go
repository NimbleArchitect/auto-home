package pluginManager

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
	"unicode"
)

type Generic struct {
	Method string
	Id     int
	Data   *json.RawMessage
}

type trigger struct {
	Name   string
	Call   string
	Fields interface{}
}

type create struct {
	Name   string
	Fields map[string]interface{}
}

type result struct {
	Ok      bool
	Message string                 `json:",omitempty"`
	Data    map[string]interface{} `json:",omitempty"`
}

type PluginConnector struct {
	name         string
	c            net.Conn
	lock         sync.Mutex
	nextId       int
	responseWait map[int]*chan result
	plug         *Plugin
	funcList     map[string]*Caller // list of function names provided by the plugin
	jsCallBack   func(string, string, map[string]interface{})
	wg           *sync.WaitGroup
}

func (c *PluginConnector) Name() string {
	c.lock.Lock()
	out := c.name
	c.lock.Unlock()
	return out
}

func (c *PluginConnector) All() map[string]*Caller {
	c.lock.Lock()
	out := c.funcList
	c.lock.Unlock()
	return out
}

func (c *PluginConnector) WaitAdd() int {
	i := c.nextId
	c.nextId++

	wait := make(chan result, 1)
	c.lock.Lock()
	c.responseWait[i] = &wait
	c.lock.Unlock()

	return i
}

// WaitDone used to signal that we should not wait any more
func (c *PluginConnector) WaitDone(i int, msg *string, data map[string]interface{}) {
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

func (c *PluginConnector) WaitOn(i int) (string, map[string]interface{}, bool) {
	c.lock.Lock()
	wait := *c.responseWait[i]
	c.lock.Unlock()

	out := <-wait

	return out.Message, out.Data, out.Ok
}

func (c *PluginConnector) writeB(b []byte) {
	fmt.Println("Mgr ->> sending", string(b))
	c.c.Write(b)
	c.lock.Lock()
	c.c.Write([]byte("\n\n"))
	c.lock.Unlock()
}

func (c *PluginConnector) handle() {
	var buf []byte

	defer c.c.Close()
	chError := make(chan error)

	go func() {

		for {
			tmp := make([]byte, 256)
			n, err := c.c.Read(tmp)
			if err != nil {
				chError <- err
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
							fmt.Println("Mgr -<< recieved", string(buf[last:i-2]))
							go c.decode(buf[last:i])
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
		case err := <-chError:
			if io.EOF == err {
				if len(buf) > 0 {
					c.decode(buf)
				}
				log.Println("connection dropped message", err)
				return
			}
		case <-time.After(1 * time.Minute):
			c.writeB([]byte{})
		}
	}
}

func (c *PluginConnector) decode(buf []byte) {
	var generic Generic

	err := json.Unmarshal(buf, &generic)
	if err != nil {
		log.Println("decode error", err)
		resp := makeError(generic.Id, err)
		c.writeB(resp)
	} else {
		//process message
		c.processMessage(generic)
	}

}

func (c *PluginConnector) processMessage(obj Generic) error {
	// runtime := goja.New()
	// runtime.SetFieldNameMapper(goja.UncapFieldNameMapper())

	// var jsobj *goja.Object
	raw, _ := obj.Data.MarshalJSON()

	switch obj.Method {
	case "trigger":
		var m trigger
		if err := json.Unmarshal(raw, &m); err != nil {
			return err
		}

		switch field := m.Fields.(type) {
		case map[string]interface{}:
			// TODO: is is worth capturing the return vars from the js call?
			// run the js callback caller.Run
			c.jsCallBack(m.Name, m.Call, field)
		default:
			errMsg := fmt.Sprintf("I don't know about type %T!\n", field)
			log.Println("trigger error:", errMsg)
			return errors.New(errMsg)
		}

		out := makeResponse(obj.Id, nil)
		c.WaitDone(obj.Id, nil, nil)
		c.writeB(out)

	case "result":
		var m result
		err := json.Unmarshal(raw, &m)
		if err != nil {
			fmt.Println("err:", err)
		}

		// proceess the response from the client trigger?
		if m.Ok {
			c.WaitDone(obj.Id, nil, m.Data)
		} else {
			c.WaitDone(obj.Id, &m.Message, nil)
		}

	case "create":
		var m create
		if err := json.Unmarshal(raw, &m); err != nil {
			if c.wg != nil {
				c.wg.Done()
			}
			return err
		}

		c.name = m.Name
		for rawName := range m.Fields {
			tmpName := []rune(rawName)
			tmpName[0] = unicode.ToLower(tmpName[0])
			name := string(tmpName)
			caller := Caller{
				Name: m.Name,
				Call: name,
				c:    c,
			}
			c.funcList[name] = &caller
		}

		out := makeResponse(obj.Id, nil)
		if c.wg != nil {
			c.wg.Done()
		}
		c.WaitDone(obj.Id, nil, nil)
		c.writeB(out)
	}

	return nil
}
