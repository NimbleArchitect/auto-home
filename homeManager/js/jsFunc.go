package js

import (
	"fmt"
	"net/rpc"
)

func (r *JavascriptVM) NewPlugin(name string, vals *rpc.Client) {
	fmt.Println(">> add plugin", name, vals)
	r.pluginList[name] = vals

}
