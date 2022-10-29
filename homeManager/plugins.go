package home

import (
	"log"
	"os"
	"os/exec"
	"path"
	"server/homeManager/pluginManager"
	"sync"
)

const SockAddr = "/tmp/rpc.sock"

type Result struct {
	Ok   bool
	Data map[string]interface{}
}

func (m *Manager) startPlugin(pluginName string, wg *sync.WaitGroup) {
	// var pluginsStarted int

	log.Println("starting plugin", pluginName)

	pluginExec := path.Join(m.pluginPath, pluginName, pluginName)

	cmd := exec.Command(pluginExec)
	// cmd.Dir = path.Join(m.pluginPath, pluginName)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		wg.Done()
		log.Println("startPlugin error:", err)
	} else {
		cmd.Wait()
		log.Printf("just finished %s (%d)\n", pluginName, cmd.Process.Pid)
	}

}

// StartPlugins starts the plugin manager and all the named plugins
func (m *Manager) StartPlugins(plug *pluginManager.Plugin) {
	var pluginList []string
	pluginList = append(pluginList, "telegram", "solar")

	wg := sync.WaitGroup{}

	wg.Add(1)

	// TODO: plugins/manager is better but still not happy with it
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

// callPluginObject is the call back function for when a plugin wants to fire an event
func (m *Manager) callPluginObject(pluginName string, call string, obj map[string]interface{}) {
	log.Println("plugin triggered", pluginName, call)
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
