package home

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"server/homeManager/pluginManager"
	log "server/logger"
	"sync"
)

const SockAddr = "/tmp/rpc.sock"

type Result struct {
	Ok   bool
	Data map[string]interface{}
}

func (m *Manager) startPlugin(pluginName string, wg *sync.WaitGroup) {

	// var pluginsStarted int

	log.Info("starting plugin", pluginName)

	pluginExec := path.Join(m.pluginPath, pluginName, pluginName)

	cmd := exec.Command(pluginExec)
	// cmd.Dir = path.Join(m.pluginPath, pluginName)
	cmd.Env = append(os.Environ(), fmt.Sprint("autohome_sockaddr=", SockAddr))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		wg.Done()
		log.Error("startPlugin error:", err)
	} else {
		cmd.Wait()
		log.Infof("just finished %s (%d)\n", pluginName, cmd.Process.Pid)
	}

}

// StartPlugins starts the plugin manager and all the named plugins
func (m *Manager) StartPlugins(plug *pluginManager.Plugin) {
	var pluginList []string

	// TODO: need to dynamically build this list of plugins
	// pluginList = append(pluginList, "telegram", "solar", "calendar")
	pluginList = append(pluginList, "telegram", "solar")

	wg := sync.WaitGroup{}

	wg.Add(1)

	// TODO: plugins/manager is better but still not happy with it
	pluginMgr := pluginManager.Manager{
		SockAddr:  "/tmp/rpc.sock",
		Plugins:   plug,
		WaitGroup: &wg,
	}
	go pluginMgr.Start(m.callPluginObject)

	for i := 0; i < len(pluginList); i++ {
		wg.Add(1)
		go m.startPlugin(pluginList[i], &wg)
	}

	log.Info("plugins started")
	wg.Done()
	wg.Wait()

	m.plugins = plug
	m.chStartupComplete <- true
	pluginMgr.WaitGroup = nil
	log.Info("startup complete")
}

// callPluginObject is the call back function for when a plugin wants to fire an event
func (m *Manager) callPluginObject(pluginName string, call string, obj map[string]interface{}) {

	log.Info("plugin triggered", pluginName, call)
	//TODO: call client on trigger, need to work out the client script to run

	// if vm := m.actions[deviceid].jsvm; vm == nil {

	//get next avaliable vm
	vm, id := m.GetNextVM()
	// once we have finished we make sure to add the vm id back to the channel list ready for next use
	defer m.PushVMID(id)

	if vm == nil {
		log.Error("invalid javascript vm")
	} else {
		// TODO: somewhere I need to validate the properties so I only save valid states
		log.Debug("state:", m.devices)

		// save the current state of all devices
		vm.SaveDeviceState(m.devices)
		// groups are set during vm init

		// register plugins
		// for n, v := range m.plugins.All() {
		// 	vm.NewPlugin(n, v)
		// }

		vm.RunJSPlugin(pluginName, call, obj)
	}

	log.Info("event finished")
}
