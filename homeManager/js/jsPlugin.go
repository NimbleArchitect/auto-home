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
				out := localCall.RunMultiArgs(values)
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
