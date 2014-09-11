package main

import (
	"logger"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if ok := logger.Init("dspLog.log"); !ok {
		panic("logger init error.")
	}
}
