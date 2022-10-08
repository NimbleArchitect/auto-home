package home

import "log"

const SockAddr = "/tmp/rpc.sock"

type Result struct {
	Ok   bool
	Data map[string]interface{}
}

func (m *Manager) startPluginManager() {
	// fmt.Println("starting plungin manager")
	// if err := os.RemoveAll(SockAddr); err != nil {
	// 	log.Fatal(err)
	// }

	// l, e := net.Listen("unix", SockAddr)
	// if e != nil {
	// 	log.Fatal("listen error:", e)
	// }

	// incoming, _ := l.Accept()

	// client := rpc.NewClient(incoming)

	// args := make(map[string]interface{})
	// // args["msg"] = "heelo world"

	// //this will store returned result
	// var result Result

	// client.Call("Client.RoleCall", args, &result)

	// if result.Ok {
	// 	m.NewPlugin(*client, result.Data)
	// }
}

func (m *Manager) startAllPlugins() {
	log.Println("starting plungins")

}
