package pluginManager

import (
	"encoding/json"
	"errors"
	"fmt"

	"io"
	"net"
	log "server/logger"
	"strings"
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

func (c *PluginConnector) Get(name string) *Caller {
	c.lock.Lock()
	out, ok := c.funcList[name]
	c.lock.Unlock()

	if ok {
		return out
	}
	return nil
}

func (c *PluginConnector) WaitAdd() int {
	i := c.nextId
	c.nextId++

	wait := make(chan result, 1)
	log.Debug("lock")
	c.lock.Lock()
	c.responseWait[i] = &wait
	c.lock.Unlock()
	log.Debug("unlock")

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
	log.Debug("lock")
	c.lock.Lock()
	wait := *c.responseWait[i]
	c.lock.Unlock()
	log.Debug("unlock")

	out := <-wait
	log.Debug("wait done")
	return out.Message, out.Data, out.Ok
}

// CloseInactive terminates all open requests
func (c *PluginConnector) CloseInactive() {
	c.lock.Lock()
	for i := 0; i < len(c.responseWait); i++ {
		wait, ok := c.responseWait[i]
		if ok {
			*wait <- result{Ok: false, Message: "plugin error", Data: nil}
		}
	}
	c.lock.Unlock()
}

func (c *PluginConnector) writeB(b []byte) {
	//fmt.Println("Mgr ->> sending", string(b))
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
		log.Debug("start reciever")
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
							log.Debug("Mgr -<< recieved", string(buf[last:i-2]))
							go c.decode(buf[last:i])
							last = i
						}

					}
					buf = []byte{}
				}
			}
		}
		log.Debug("finish reciever")
	}()

	for {
		select {
		case err := <-chError:
			if io.EOF == err {
				if len(buf) > 0 {
					c.decode(buf)
				}
				log.Error(c.name, "dropped connection", err)
				c.CloseInactive()
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
		log.Error("decode error", err)
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
		case nil:
			// the recieved field is empty so run the jsCallBack with empty arguments
			c.jsCallBack(m.Name, m.Call, make(map[string]interface{}))
		case map[string]interface{}:
			// TODO: is is worth capturing the return vars from the js call?
			// run the js callback caller.Run
			c.jsCallBack(m.Name, m.Call, field)
		default:
			errMsg := fmt.Sprintf("I don't know about type %T!\n", field)
			log.Error("trigger error:", errMsg)
			return errors.New(errMsg)
		}

		out := makeResponse(obj.Id, nil)
		c.WaitDone(obj.Id, nil, nil)
		c.writeB(out)

	case "result":
		var m result
		err := json.Unmarshal(raw, &m)
		if err != nil {
			log.Error("err:", err)
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

		var callList []string
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
			log.Debug("add function:", caller)
			c.funcList[name] = &caller
			callList = append(callList, name)
		}

		out := makeResponse(obj.Id, nil)
		if c.wg != nil {
			c.wg.Done()
		}
		c.WaitDone(obj.Id, nil, nil)
		c.writeB(out)
		log.Infof("plugin \"%s\" registered: %s\n", m.Name, strings.Join(callList, ", "))
	}

	return nil
}
