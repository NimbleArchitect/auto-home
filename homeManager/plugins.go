package home

import (
	"log"
	"os"
	"os/exec"
	"path"
	"server/homeManager/pluginManager"
	"sync"

	"github.com/dop251/goja"
)

const SockAddr = "/tmp/rpc.sock"

type Result struct {
	Ok   bool
	Data map[string]interface{}
}

// func (m *Manager) startPluginManager(done chan int) {
// 	// TODO: the whole plugin system is a bit crap... needs rewriting
// 	var pluginsRegistered int

// 	fmt.Println("starting plugin manager")
// 	if err := os.RemoveAll(SockAddr); err != nil {
// 		log.Fatal(err)
// 	}

// 	l, e := net.Listen("unix", SockAddr)
// 	if e != nil {
// 		log.Fatal("listen error:", e)
// 	}

// 	done <- 0

// 	for {
// 		incoming, _ := l.Accept()
// 		go func() {
// 			// fmt.Println(">> recieved plugin connection")
// 			client := rpc.NewClient(incoming)

// 			args := make(map[string]interface{})

// 			//this will store returned result
// 			var result Result

// 			client.Call("Client.RoleCall", args, &result)

// 			if result.Ok {
// 				name := result.Data["name"].(string)
// 				log.Println("allowing plugin", name)
// 				m.plugins[name] = client
// 				pluginsRegistered++
// 				done <- pluginsRegistered
// 			}
// 		}()
// 	}
// }

func (m *Manager) startPlugin(pluginName string, wg *sync.WaitGroup) {
	// var pluginsStarted int

	log.Println("starting plungin", pluginName)

	pluginExec := path.Join(m.pluginPath, pluginName, pluginName)

	cmd := exec.Command(pluginExec)
	// cmd.Dir = path.Join(m.pluginPath, pluginName)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
		wg.Done()
	}

	cmd.Wait()
	log.Printf("just finished %s (%d)\n", pluginName, cmd.Process.Pid)
}

// TODO:: plugins dont work needs work
func (m *Manager) StartPlugins(plug *pluginManager.Plugin) {
	var pluginList []string
	pluginList = append(pluginList, "telegram")

	wg := sync.WaitGroup{}

	wg.Add(1)

	// TODO: plugins/manager need a rewrite
	go pluginManager.Start("/tmp/rpc.sock", &wg, plug, m.callPluginObject)

	for i := 0; i < len(pluginList); i++ {
		wg.Add(1)
		go m.startPlugin(pluginList[i], &wg)
	}

	log.Println("plugins started")
	wg.Done()
	wg.Wait()

	m.plugins = plug
	m.chStartupComplete <- true

	log.Println("startup complete")
}

func (m *Manager) callPluginObject(pluginName string, call string, obj *goja.Object) {
	log.Println("plugin triggered")
	//TODO: call client on trigger, need to work out the client script to run

	// if vm := m.actions[deviceid].jsvm; vm == nil {

	//get next avaliable vm
	vm, id := m.GetNextVM()
	// once we have finished we make sure to add the vm id back to the channel list ready for next use
	defer m.PushVMID(id)

	if vm == nil {
		log.Println("invalid javascript vm")
	} else {
		// TODO: somewhere I need to validate the properties so I only save valid states
		log.Println("state:", m.devices)

		// save the current state of all devices
		vm.SaveDeviceState(m.devices)
		// groups are set during vm init

		// register plugins
		// for n, v := range m.plugins.All() {
		// 	vm.NewPlugin(n, v)
		// }

		vm.RunJSPlugin(pluginName, call, obj)
	}

	log.Println("event finished")
}
