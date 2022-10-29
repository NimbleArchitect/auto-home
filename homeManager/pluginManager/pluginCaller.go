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

	nextId := c.c.WaitAdd()

	fields := make(map[int]interface{})

	log.Printf("called: %s.%s(%s)\n", c.Name, c.Call, values)
	for i, v := range values {
		fields[i] = v.Export()
	}
	t.Name = c.Name
	t.Call = c.Call
	t.Fields = fields

	generic := response{
		Method: "trigger",
		Id:     nextId,
		Data:   t,
	}
	data, err := json.Marshal(generic)
	if err != nil {
		log.Println("json error", err)
	}

	c.c.writeB(data)
	c.c.WaitOn(nextId)
}
