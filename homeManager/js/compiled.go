package js

import (
	"log"
	"os"
	"strings"

	"github.com/dop251/goja"
)

type CompiledScripts struct {
	compiled map[string]*goja.Program
}

func loadScript(filename string) *goja.Program {
	cfile, err := os.ReadFile(filename)
	if err != nil {
		log.Println("unable to read file:", err)
	}

	prog, err := goja.Compile(filename, ";(function () {"+string(cfile)+"\n})", true)

	// prog, err := goja.Compile(filename, string(cfile), true)
	if err != nil {
		log.Println("unable to compile script", err)
	}

	return prog
}

func (c *CompiledScripts) NewVM() (*JavascriptVM, error) {
	var console jsConsole

	runtime := goja.New()
	runtime.SetFieldNameMapper(goja.UncapFieldNameMapper())

	vm := JavascriptVM{
		runtime:     runtime,
		deviceCode:  make(map[string]*goja.Object),
		deviceState: make(map[string]jsDevice),
		groupCode:   make(map[string]*goja.Object),
		groups:      make(map[string]jsGroup),
	}

	err := runtime.Set("console", console)
	if err != nil {
		return nil, err
	}

	err = runtime.Set("thread", runAsThread)
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
				log.Println("internal: not a function")
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
			// fmt.Println(p)
			if p != nil {
				// log.Println("saving compiled code")
				compiled[item.Name()] = p
			}
		}
	}

	return CompiledScripts{
		compiled: compiled,
	}
}
