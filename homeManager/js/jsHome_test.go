package js

import (
	"server/globals"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/dop251/goja"
)

func setupVM() JavascriptVM {
	g := globals.New()

	return JavascriptVM{
		waitGroup: sync.Mutex{},
		global:    g,
	}
}

func TestGetDeviceByName(t *testing.T) {
	vm := setupVM()

	home := jsHome{
		vm:      &vm,
		devices: make(map[string]jsDevice),
	}

	home.devices["dev1"] = jsDevice{Name: "name1"}
	home.devices["dev2"] = jsDevice{Name: "name2"}
	home.devices["dev3"] = jsDevice{Name: "name3"}

	dev := home.GetDeviceByName("name2")

	if dev.Name != "name2" {
		t.Fatalf("expected name2 got %s", dev.Name)
	}
}

func TestGetDeviceById(t *testing.T) {
	vm := setupVM()

	home := jsHome{
		vm:      &vm,
		devices: make(map[string]jsDevice),
	}

	home.devices["dev1"] = jsDevice{Name: "name1", Id: "dev1"}
	home.devices["dev2"] = jsDevice{Name: "name2", Id: "dev2"}
	home.devices["dev3"] = jsDevice{Name: "name3", Id: "dev3"}

	dev := home.GetDeviceById("dev2")

	if dev.Id != "dev2" {
		t.Fatalf("expected dev2 got %s", dev.Id)
	}
}

func TestGetDevices(t *testing.T) {
	vm := setupVM()

	home := jsHome{
		vm:      &vm,
		devices: make(map[string]jsDevice),
	}

	home.devices["dev1"] = jsDevice{Name: "name1", Id: "dev1"}
	home.devices["dev2"] = jsDevice{Name: "name2", Id: "dev2"}
	home.devices["dev3"] = jsDevice{Name: "name3", Id: "dev3"}

	dev := home.GetDevices()

	if len(dev) != 3 {
		t.Fatalf("expected 3 got %d", len(dev))
	}

	arrNames := make([]string, 3)
	for _, v := range dev {
		if v.Name == "name1" {
			arrNames[0] = v.Name
		}
		if v.Name == "name2" {
			arrNames[1] = v.Name
		}
		if v.Name == "name3" {
			arrNames[2] = v.Name
		}
	}
	found := strings.Join(arrNames, " ")

	if found != "name1 name2 name3" {
		t.Fatalf("expected \"name1 name2 name3\" got \"%s\"", found)
	}
}

func TestGetDevicesStartName(t *testing.T) {
	vm := setupVM()

	home := jsHome{
		vm:      &vm,
		devices: make(map[string]jsDevice),
	}

	home.devices["dev1"] = jsDevice{Name: "name1", Id: "dev1"}
	home.devices["dev2"] = jsDevice{Name: "notname2", Id: "dev2"}
	home.devices["dev3"] = jsDevice{Name: "name3", Id: "dev3"}

	dev := home.GetDevicesStartName("not")

	if len(dev) != 1 {
		t.Fatalf("expected 3 got %d", len(dev))
	}

	var arrNames []string
	for _, v := range dev {
		arrNames = append(arrNames, v.Name)
	}
	found := strings.Join(arrNames, " ")

	if found != "notname2" {
		t.Fatalf("expected \"notname2\" got \"%s\"", found)
	}
}

func TestCountdownRestart(t *testing.T) {
	const delay = 0.0100
	delayDuration := (delay * 1000) * time.Millisecond

	vm := setupVM()

	home := jsHome{
		vm: &vm,
	}

	for i := 0; i < 10; i++ {
		time.Sleep(25 * time.Millisecond)
		home.Countdown("trial", delay, goja.Undefined())
	}

	start := time.Now()
	vm.Wait()
	duration := time.Since(start)

	if duration < delayDuration {
		t.Fatalf("Sleep(%s) slept for only %s", delayDuration, duration)
	}
}

func TestCountdownCancel(t *testing.T) {
	const delay = 0.0100
	delayDuration := (delay * 1000) * time.Millisecond

	vm := setupVM()

	home := jsHome{
		vm: &vm,
	}

	for i := 0; i < 10; i++ {
		time.Sleep(25 * time.Millisecond)
		home.Countdown("trial", 0.0100, goja.Undefined())
	}

	home.Countdown("trial", 0, goja.Undefined())
	start := time.Now()
	vm.Wait()
	duration := time.Since(start)

	if duration > delayDuration {
		t.Fatalf("Sleep(%s) slept for too long %s", delayDuration, duration)
	}
}
