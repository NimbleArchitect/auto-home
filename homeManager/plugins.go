package home

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"time"
)

const SockAddr = "/tmp/rpc.sock"

type Result struct {
	Ok   bool
	Data map[string]interface{}
}

func (m *Manager) startPluginManager(done chan bool) {
	fmt.Println("starting plugin manager")
	if err := os.RemoveAll(SockAddr); err != nil {
		log.Fatal(err)
	}

	l, e := net.Listen("unix", SockAddr)
	if e != nil {
		log.Fatal("listen error:", e)
	}

	done <- true

	for {
		incoming, _ := l.Accept()
		go func() {
			// fmt.Println(">> recieved plugin connection")
			client := rpc.NewClient(incoming)

			args := make(map[string]interface{})

			//this will store returned result
			var result Result

			client.Call("Client.RoleCall", args, &result)

			if result.Ok {
				name := result.Data["name"].(string)
				log.Println("allowing plugin", name)
				m.plugins[name] = client

			}
		}()
	}
}

func (m *Manager) startAllPlugins(done chan bool) {
	log.Println("starting plungins")

	cmd := exec.Command("plugins/telegram/telegram")
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	done <- true
	cmd.Wait()
	log.Printf("Just ran subprocess %d, exiting\n", cmd.Process.Pid)
}

// function (m *Manager) CallPlugin(pluginname string, args string) {
// 	log.Println("starting plungins")

// }
