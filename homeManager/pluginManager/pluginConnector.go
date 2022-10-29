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
	Message string
}

type PluginConnector struct {
	name       string
	c          net.Conn
	lock       sync.Mutex
	nextId     int
	wait       map[int]chan bool
	plug       *Plugin
	funcList   map[string]*Caller
	jsCallBack func(string, string, map[string]interface{})
	wg         *sync.WaitGroup
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
			c.c.SetWriteDeadline(time.Now().Add(5 * time.Second))
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
		c.c.SetWriteDeadline(time.Now().Add(5 * time.Second))
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

		// TODO: trigger dosent work, needs partially moving to js VM and we need to use c.jsCallBack
		switch field := m.Fields.(type) {
		case map[string]interface{}:
			c.jsCallBack(m.Name, m.Call, field)
		default:
			errMsg := fmt.Sprintf("I don't know about type %T!\n", field)
			log.Println("trigger error:", errMsg)
			return errors.New(errMsg)
		}

		out := makeError(obj.Id, nil)
		c.writeB(out)
		c.lock.Lock()
		c.wait[obj.Id] <- true
		c.lock.Unlock()

	case "result":
		var m result
		json.Unmarshal(raw, &m)
		c.lock.Lock()
		c.wait[obj.Id] <- true
		c.lock.Unlock()

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
			name := []rune(rawName)
			name[0] = unicode.ToLower(name[0])
			caller := Caller{
				Name: m.Name,
				Call: string(name),
				c:    c,
			}

			c.funcList[string(name)] = &caller
		}

		out := makeError(obj.Id, nil)
		if c.wg != nil {
			c.wg.Done()
		}
		c.writeB(out)
		c.lock.Lock()
		c.wait[obj.Id] <- true
		c.lock.Unlock()
	}

	return nil
}
