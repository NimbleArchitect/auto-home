package home

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
)

const SockAddr = "/tmp/rpc.sock"

type Result struct {
	Ok   bool
	Data map[string]interface{}
}

func (m *Manager) startPluginManager() {
	fmt.Println("starting plungin manager")
	if err := os.RemoveAll(SockAddr); err != nil {
		log.Fatal(err)
	}

	l, e := net.Listen("unix", SockAddr)
	if e != nil {
		log.Fatal("listen error:", e)
	}

	for {
		incoming, _ := l.Accept()
		// fmt.Println(">> recieved plugin connection")
		client := rpc.NewClient(incoming)

		args := make(map[string]interface{})
		// args["msg"] = "heelo world"

		//this will store returned result
		var result Result

		client.Call("Client.RoleCall", args, &result)

		if result.Ok {
			name := result.Data["name"].(string)
			log.Println("allowing plugin", name)
			m.plugins[name] = client

		}
	}
}

func (m *Manager) startAllPlugins() {
	log.Println("starting plungins")

}

// function (m *Manager) CallPlugin(pluginname string, args string) {
// 	log.Println("starting plungins")

// }
