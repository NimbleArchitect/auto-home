package js

import (
	"log"
	"net/rpc"

	"github.com/dop251/goja"
)

type jsPlugin struct {
	client *rpc.Client
	name   string
}

type Result struct {
	Ok   bool
	Data map[string]interface{}
}

func (r *JavascriptVM) NewPlugin(name string, vals *rpc.Client) {
	fmt.Println(">> add plugin", name, vals)
	r.pluginList[name] = vals

}

func (d *jsPlugin) Call(funcName string, vars goja.Value) {
	var result Result

	if d.client != nil {
		args := make(map[string]interface{})
		log.Println("calling", d.name+"."+funcName)

		keys := vars.(*goja.Object).Keys()
		for _, name := range keys {
			v := vars.(*goja.Object).Get(name).String()
			args[name] = v
		}

		err := d.client.Call(d.name+"."+funcName, args, &result)
		if err != nil {
			log.Println("error calling", funcName, err)
		}

		if result.Ok {
			log.Println("plugin returned ok")
		}
		log.Println("plugin response:", result.Data)

	}
}
