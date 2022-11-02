package js

import (
	"github.com/dop251/goja"
)

// type jsPlugin struct {
// 	client *rpc.Client
// 	name   string
// }

type Result struct {
	Ok   bool
	Data map[string]interface{}
}

func (r *JavascriptVM) loadPlugins() {
	for n, plugin := range r.pluginList.All() {
		thisPlugin := r.runtime.NewObject()
		for name, caller := range plugin.All() {
			// de-reference caller so we get a copy of it that we can use in the function
			localCall := *caller

			thisPlugin.Set(name, func(values ...goja.Value) goja.Value {
				out := localCall.Run(values)
				if len(out) == 0 {
					return goja.Undefined()
				}
				if len(out) == 1 {
					for _, v := range out {
						return r.runtime.ToValue(v)
					}
				}
				return r.runtime.ToValue(out)
			})
		}
		r.plugins[n] = thisPlugin
	}
}

// func (r *JavascriptVM) NewPlugin(name string, vals *pluginManager.PluginConnector) {
// 	r.pluginList[name] = vals
// }

// func (d *jsPlugin) Call(funcName string, vars goja.Value) {
// 	var result Result
// 	if d.client != nil {
// 		args := make(map[string]interface{})
// 		log.Println("calling", d.name+"."+funcName)

// 		keys := vars.(*goja.Object).Keys()
// 		for _, name := range keys {
// 			v := vars.(*goja.Object).Get(name).String()
// 			args[name] = v
// 		}

// 		err := d.client.Call(d.name+"."+funcName, args, &result)
// 		if err != nil {
// 			log.Println("error calling", funcName, err)
// 		}

// 		if result.Ok {
// 			log.Println("plugin returned ok")
// 		}
// 		log.Println("plugin response:", result.Data)

// 	}
// }
