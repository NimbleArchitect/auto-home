package js

import (
	"log"
	"os"
	"server/homeManager/pluginManager"
	"strings"

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
		log.Println("unable to read file:", err)
	}

	if strings.HasSuffix(filename, "/common.js") {
		prog, err = goja.Compile(filename, string(cfile), true)
	} else {
		prog, err = goja.Compile(filename, ";(function () {"+string(cfile)+"\n})", true)
	}

	// prog, err := goja.Compile(filename, string(cfile), true)
	if err != nil {
		log.Println("unable to compile script", err)
	}

	return prog
}

// func (c *CompiledScripts) NewVM(pluginList *pluginManager.Plugin) (*JavascriptVM, error) {
func (c *CompiledScripts) NewVM(pluginList *pluginManager.Plugin) (*JavascriptVM, error) {
	var console jsConsole

	runtime := goja.New()
	runtime.SetFieldNameMapper(goja.UncapFieldNameMapper())

	vm := JavascriptVM{
		runtime:     runtime,
		deviceCode:  make(map[string]*goja.Object),
		deviceState: make(map[string]jsDevice),
		groupCode:   make(map[string]*goja.Object),
		groups:      make(map[string]jsGroup),
		pluginList:  pluginList,
		pluginCode:  make(map[string]*goja.Object),
	}

	plugins := runtime.NewObject()
	for n, plugin := range pluginList.All() {
		thisPlugin := runtime.NewObject()
		for _, caller := range plugin.All() {

			// thisPlugin.Set(caller.Call, caller.Run)
			thisPlugin.Set(caller.Call, func(values ...goja.Value) goja.Value {
				out := caller.Run(values)
				if len(out) == 0 {
					return goja.Undefined()
				}
				if len(out) == 1 {
					for _, v := range out {
						return vm.runtime.ToValue(v)
					}
				}
				return vm.runtime.ToValue(out)
			})
		}

		plugins.Set(n, thisPlugin)
	}
	//TODO: move this to runJS, to do that I need to move plugins to a pointer that is attached to the VM
	runtime.Set("plugin", plugins)

	err := runtime.Set("console", console)
	if err != nil {
		return nil, err
	}

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
			log.Println("error running script", scriptName, err)
		} else {
			call, ok := goja.AssertFunction(module)
			if ok {
				_, err := call(goja.Undefined())
				if err != nil {
					log.Println("script error", err)
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
			log.Println("loading script", fullname)
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
