package main

import (
	mon "Zetsu/monitor/func"
	"time"
)

const (
	address = "localhost:5790"
)

func main() {
	monitor := mon.NewMonitor(address)
	time.Sleep(time.Second * 20)
	monitor.Stop()
}
