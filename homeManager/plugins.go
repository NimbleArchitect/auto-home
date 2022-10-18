package home

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/exec"
)

const SockAddr = "/tmp/rpc.sock"

type Result struct {
	Ok   bool
	Data map[string]interface{}
}

func (m *Manager) startPluginManager(done chan int) {
	// TODO: the whole plugin system is a bit crap... needs rewriting
	var pluginsRegistered int

	fmt.Println("starting plugin manager")
	if err := os.RemoveAll(SockAddr); err != nil {
		log.Fatal(err)
	}

	l, e := net.Listen("unix", SockAddr)
	if e != nil {
		log.Fatal("listen error:", e)
	}

	done <- 0

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
				pluginsRegistered++
				done <- pluginsRegistered
			}
		}()
	}
}

func (m *Manager) startPlugin(pluginName string) {
	// var pluginsStarted int

	log.Println("starting plungin", pluginName)

	cmd := exec.Command(m.pluginPath + pluginName + "/" + pluginName)
	cmd.Dir = m.pluginPath + pluginName + "/"

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	cmd.Wait()
	log.Printf("just finished %s (%d)\n", pluginName, cmd.Process.Pid)
}

// function (m *Manager) CallPlugin(pluginname string, args string) {
// 	log.Println("starting plungins")

// }
