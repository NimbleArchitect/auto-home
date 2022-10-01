package js

import (
	"fmt"
	"strings"
	"time"
)

type jsConsole struct {
	start time.Time
}

func (d *jsConsole) Log(s ...string) {
	fmt.Println(strings.Join(s, " "))
}

func (d *jsConsole) StartTimer() {
	d.start = time.Now()
}

func (d *jsConsole) StopTimer() {
	elapsed := time.Since(d.start)
	fmt.Printf("\n\ntime taken %s\n\n", elapsed)
}
