package js

import (
	"os"
	"server/globals"
	"server/homeManager/pluginManager"
	log "server/logger"
	"strings"
	"sync"

	"github.com/dop251/goja"
)

type CompiledScripts struct {
	compiled map[string]*goja.Program
}

func loadScript(filename string) *goja.Program {
	var prog *goja.Program
	var err error

	cfile, err := os.ReadFile(filename)
	if err != nil {
		log.Error("unable to read file:", err)
	}

	if strings.HasSuffix(filename, "/common.js") {
		prog, err = goja.Compile(filename, string(cfile), true)
	} else {
		prog, err = goja.Compile(filename, ";(function () {"+string(cfile)+"\n})", true)
	}

	// prog, err := goja.Compile(filename, string(cfile), true)
	if err != nil {
		log.Error("unable to compile script", err)
	}

	return prog
}

// func (c *CompiledScripts) NewVM(pluginList *pluginManager.Plugin) (*JavascriptVM, error) {
func (c *CompiledScripts) NewVM(pluginList *pluginManager.Plugin, global *globals.Global) (*JavascriptVM, error) {
	var console jsConsole

	runtime := goja.New()
	runtime.SetFieldNameMapper(goja.UncapFieldNameMapper())

	waitGroup := sync.WaitGroup{}

	vm := JavascriptVM{
		runtime:     runtime,
		global:      global,
		deviceCode:  make(map[string]*goja.Object),
		deviceState: make(map[string]jsDevice),
		groupCode:   make(map[string]*goja.Object),
		groups:      make(map[string]jsGroup),
		pluginList:  pluginList,
		pluginCode:  make(map[string]*goja.Object),
		plugins:     make(map[string]*goja.Object),
		vmInUseLock: &waitGroup,
	}

	vm.loadPlugins()

	err := runtime.Set("console", console)
	if err != nil {
		return nil, err
	}

	// TODO: the calling vm needs to wait for all threaded calls to finish before running to completion,
	//  doing so will mean that threaded calls dont accidently get added to the history when they are external
	err = runtime.Set("thread", vm.runAsThread)
	if err != nil {
		return nil, err
	}

	// TODO: I still need to add ability to set common functions
	err = runtime.Set("set", vm.objLoader)
	if err != nil {
		return nil, err
	}

	// load all scripts one after the other and call the
	//  returned object
	for scriptName, code := range c.compiled {
		// run the script module
		module, err := runtime.RunProgram(code)
		if err != nil {
			log.Error("error running script", scriptName, err)
		} else {
			call, ok := goja.AssertFunction(module)
			if ok {
				_, err := call(goja.Undefined())
				if err != nil {
					log.Error("script error", err)
				}
			} else {
				// TODO: should this be enabled?
				// log.Println("internal: not a function")
			}
		}
	}

	return &vm, nil
}

func LoadAllScripts(path string) CompiledScripts {

	compiled := make(map[string]*goja.Program)

	sep := string(os.PathSeparator)

	pathname := strings.TrimSuffix(path, sep)
	entires, err := os.ReadDir(pathname)
	if err != nil {
		return CompiledScripts{}
	}

	for _, item := range entires {
		if !item.IsDir() {
			fullname := pathname + sep + item.Name()
			log.Info("loading script", fullname)
			p := loadScript(fullname)
			if p != nil {
				compiled[item.Name()] = p
			}
		}
	}

	return CompiledScripts{
		compiled: compiled,
	}
}
