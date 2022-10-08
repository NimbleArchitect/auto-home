package js

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type jsConsole struct {
	start time.Time
}

func (d *jsConsole) Log(s ...string) {
	log.Println(strings.Join(s, " "))
}

func (d *jsConsole) StartTimer() {
	d.start = time.Now()
}

func (d *jsConsole) StopTimer() {
	elapsed := time.Since(d.start)
	fmt.Printf("\ntime taken %s\n", elapsed)
}
