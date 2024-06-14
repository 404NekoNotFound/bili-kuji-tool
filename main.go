package main

import (
	"bili-kuji-management/src/logger"
	"bili-kuji-management/src/tui"
)

func main() {
	log := logger.New("kuji.log", true, true)

	tui.Ready(tui.WithLogger(log)).Run()
}
