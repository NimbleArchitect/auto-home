package js

import (
	"fmt"
	"server/globals"
	log "server/logger"
	"sync"
	"testing"

	"github.com/dop251/goja"
)

func setupJsVM() JavascriptVM {
	g := globals.New()

	runtime := goja.New()
	runtime.SetFieldNameMapper(goja.UncapFieldNameMapper())

	return JavascriptVM{
		waitGroup: sync.Mutex{},
		global:    g,
		runtime:   runtime,
	}
}

func TestJavaScriptObject(t *testing.T) {
	vm := setupJsVM()

	js := &JavascriptVM{
		runtime: vm.runtime,
	}

	home := jsHome{
		vm:      &vm,
		devices: make(map[string]jsDevice),
	}

	dial := make(map[string]jsDial)
	dial["state"] = jsDial{
		Name:     "state",
		Value:    5,
		previous: 4,
		min:      1,
		max:      12,
		flag:     jsFlag{},
	}

	home.devices["test-1-id"] = jsDevice{
		js:       js,
		Name:     "test-1",
		Id:       "test-1-id",
		propDial: dial,
	}

	if err := vm.runtime.Set("home", home); err != nil {
		t.Error("unable to attach plugin object to javascript vm:", err)
	}

	out, err := vm.runtime.RunString("home.getDeviceByName('test-1').get('state').value")
	if err != nil {
		t.Fatal("error running javascript:", err)
	}

	if out.ToInteger() != 5 {
		t.Fatalf("expected 5 got %s", out.String())
	}
}

type HomeObjectResult struct {
	funcName   string
	funcScript string
	expected   interface{}
}

func TestJsHomeObject(t *testing.T) {

	table := []HomeObjectResult{
		{
			funcName:   "test_getDeviceByName_dial1_value",
			funcScript: `return home.getDeviceByName('test-1').get('dial1').value`,
			expected:   int64(5),
		}, {
			funcName:   "test_getDeviceByName_dial1_previous",
			funcScript: `return home.getDeviceByName('test-1').get('dial1').previous`,
			expected:   int64(4),
		}, {
			funcName:   "test_getDeviceById_dial1_value",
			funcScript: `return home.getDeviceById('test-1-id').get('dial1').value`,
			expected:   int64(5),
		}, {
			funcName:   "test_getDeviceById_dial1_previous",
			funcScript: `return home.getDeviceById('test-1-id').get('dial1').previous`,
			expected:   int64(4),
		}, {
			funcName:   "test_getDeviceById_switch1_value",
			funcScript: `return home.getDeviceById('test-1-id').get('switch1').value`,
			expected:   bool(true),
		}, {
			funcName:   "test_getDeviceById_switch1_previous",
			funcScript: `return home.getDeviceById('test-1-id').get('switch1').previous.valueOf()`,
			expected:   bool(false),
		},
		// TODO: needs finishing, values recieved are not correct, object prototype isnt being called.
		//  I shouldn't need to call .valueOf() manually
		{
			funcName:   "test_getDeviceById_switch2_value",
			funcScript: `return home.getDeviceById('test-1-id').get('switch2').value.valueOf()`,
			expected:   bool(false),
		}, {
			funcName:   "test_getDeviceById_switch2_previous",
			funcScript: `return home.getDeviceById('test-1-id').get('switch2').previous.valueOf()`,
			expected:   bool(true),
		},
	}

	dial := make(map[string]jsDial)
	dial["dial1"] = jsDial{
		Name:     "dial1",
		Value:    5,
		previous: 4,
		min:      1,
		max:      12,
		flag:     jsFlag{},
	}

	swi := make(map[string]jsSwitch)
	swi["switch1"] = jsSwitch{
		Name:     "switch1",
		Value:    "open",
		state:    true,
		previous: "close",
		flag:     jsFlag{},
	}
	swi["switch2"] = jsSwitch{
		Name:     "switch2",
		Value:    "close",
		state:    false,
		previous: "open",
		flag:     jsFlag{},
	}

	jsScriptCode := createSetScriptFromTable("test1", table)
	vm, err := buildVmFromScript(jsScriptCode)
	if err != nil {
		t.Error("unable to build script:", err)
	}

	vm.deviceState = make(map[string]jsDevice)
	vm.deviceState["test-1-id"] = jsDevice{
		js:         vm,
		Name:       "test-1",
		Id:         "test-1-id",
		propDial:   dial,
		propSwitch: swi,
	}

	for _, value := range table {
		out, err := vm.RunJS("test1", value.funcName, goja.Undefined())
		if err != nil {
			t.Error("unable to call function", value.funcName, ":", err)
		}

		switch v := value.expected.(type) {
		case int64:
			if out.ToInteger() != v {
				t.Fatalf("check %s expected \"%d\", recieved \"%d\"", value.funcName, v, out.ToInteger())
			}
		case string:
			if out.ToString().String() != v {
				t.Fatalf("check %s expected \"%s\", recieved \"%s\"", value.funcName, v, out.ToString().String())
			}
		case bool:
			if out.ToBoolean() != v {
				t.Fatalf("check \"%s\" expected %t, recieved %t\"", value.funcName, v, out.ToBoolean())
			}
		default:
			fmt.Printf("I don't know about type %T!\n", v)
		}

		// fmt.Println("out >>", out)
	}
}

func createSetScriptFromTable(setName string, table []HomeObjectResult) string {
	scriptFunctions := ""
	for _, v := range table {
		scriptFunctions += fmt.Sprintf(`%s(props) {
		%s
	},
	`, v.funcName, v.funcScript)
	}

	jsScriptCode := fmt.Sprintf(`set("%s", {
	%s
})`, setName, scriptFunctions[:len(scriptFunctions)-1])

	return jsScriptCode
}

func buildVmFromScript(script string) (*JavascriptVM, error) {

	prog, err := goja.Compile("test", ";(function () {"+script+"\n})", true)
	if err != nil {
		// errors.New("unable to compile javascript:", err)
		return nil, err
	}

	var console jsConsole

	runtime := goja.New()
	runtime.SetFieldNameMapper(goja.UncapFieldNameMapper())

	vm := JavascriptVM{
		runtime:     runtime,
		global:      nil,
		deviceCode:  make(map[string]*goja.Object),
		deviceState: make(map[string]jsDevice),
		groupCode:   make(map[string]*goja.Object),
		groups:      make(map[string]jsGroup),
		pluginList:  nil,
		pluginCode:  make(map[string]*goja.Object),
		plugins:     make(map[string]*goja.Object),
	}

	err = runtime.Set("console", console)
	if err != nil {
		// t.Error("unable to connect console:", err)
		return nil, err
	}

	err = runtime.Set("thread", vm.runAsThread)
	if err != nil {
		// t.Error("unable to connect thread:", err)
		return nil, err
	}

	// TODO: I still need to add ability to set common functions
	err = runtime.Set("set", vm.objLoader)
	if err != nil {
		// t.Error("unable to connect set:", err)
		return nil, err
	}

	module, err := runtime.RunProgram(prog)
	if err != nil {
		// t.Error("error running script:", err)
		return nil, err
	} else {
		call, ok := goja.AssertFunction(module)
		if ok {
			_, err := call(goja.Undefined())
			if err != nil {
				log.Error("script error", err)
			}
		}
	}

	return &vm, nil
}
