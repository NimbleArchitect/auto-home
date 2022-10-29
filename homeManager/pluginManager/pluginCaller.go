package pluginManager

import (
	"encoding/json"
	"log"

	"github.com/dop251/goja"
)

type Caller struct {
	Name string
	Call string
	c    *PluginConnector
}

func (c *Caller) Run(values ...goja.Value) {
	var t trigger

	fields := make(map[int]interface{})
	c.c.nextId++

	log.Printf("called: %s.%s()\n", c.Name, c.Call)
	for i, v := range values {
		fields[i] = v.Export()
	}
	t.Name = c.Name
	t.Call = c.Call
	t.Fields = fields

	generic := response{
		Method: "trigger",
		Id:     c.c.nextId,
		Data:   t,
	}
	data, err := json.Marshal(generic)
	if err != nil {
		log.Println("json error", err)
	}

	c.c.lock.Lock()
	_, ok := c.c.wait[c.c.nextId]
	c.c.lock.Unlock()
	if !ok {
		newChan := make(chan bool, 1)
		c.c.lock.Lock()
		c.c.wait[c.c.nextId] = newChan
		<-c.c.wait[c.c.nextId]
		c.c.writeB(data)
		c.c.lock.Unlock()
		// TODO: this needs to be wrapped in a select so we can have a timeout
		// <-newChan
	}

}
