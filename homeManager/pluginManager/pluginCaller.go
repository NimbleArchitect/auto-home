package pluginManager

import (
	"encoding/json"
	log "server/logger"

	"github.com/dop251/goja"
)

type Caller struct {
	Name string
	Call string
	c    *PluginConnector
}

// RunMultiArgs sends a call to the remote plugin and waits for a response supports an array of arguments
func (c *Caller) RunMultiArgs(values []goja.Value) map[string]interface{} {

	fields := make(map[int]interface{})

	log.Infof("called: %s.%s(%s)\n", c.Name, c.Call, values)
	for i, v := range values {
		fields[i] = v.Export()
	}

	return c.Run(fields)

}

// Run sends a call to the remote plugin and waits for a response, only supports a single object
func (c *Caller) Run(fields interface{}) map[string]interface{} {
	var t trigger

	nextId := c.c.WaitAdd()

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
		log.Error("json error", err)
	}

	c.c.writeB(data)
	msg, args, ok := c.c.WaitOn(nextId)

	if !ok {
		// TODO: need to check if the plugin failed due to plugin argument problems of if the plugin crashed
		//  if the plugin crashed then if should be reloaded so we should re-queue the request, if it
		//  fails a second time then we need to skip over and return a failure
		log.Error(c.Name, "error:", msg)
	}

	log.Debug("return")
	return args

}
