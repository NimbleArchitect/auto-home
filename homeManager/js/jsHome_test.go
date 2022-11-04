// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package js

import (
	"server/globals"
	"sync"
	"testing"
	"time"

	"github.com/dop251/goja"
)

func TestCountdownRestart(t *testing.T) {
	const delay = 100
	delayDuration := delay * time.Millisecond

	g := globals.New()

	vm := JavascriptVM{
		waitGroup: sync.Mutex{},
		global:    g,
	}

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
	const delay = 500
	delayDuration := delay * time.Nanosecond

	g := globals.New()

	vm := JavascriptVM{
		waitGroup: sync.Mutex{},
		global:    g,
	}

	home := jsHome{
		vm: &vm,
	}

	for i := 0; i < 10; i++ {
		time.Sleep(25 * time.Millisecond)
		home.Countdown("trial", 100, goja.Undefined())
	}

	home.Countdown("trial", 0, goja.Undefined())
	start := time.Now()
	vm.Wait()
	duration := time.Since(start)

	if duration > delayDuration {
		t.Fatalf("Sleep(%s) slept for too long %s", delayDuration, duration)
	}
}
