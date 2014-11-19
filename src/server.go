package main

import (
	"adaptor/mango"
	"adaptor/megaMedia"
	"filter"
	"logger"
	"manager"

	//"fmt"
	"net/http"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//runtime.GOMAXPROCS(1)

	logger.WinLog = logger.NewLog("win.log")
	logger.ClickLog = logger.NewLog("click.log")
	logger.DisplayLog = logger.NewLog("display.log")
	logger.BidLog = logger.NewLog("bid.log")
	logger.GiveupLog = logger.NewLog("giveup.log")

	mango.Init("mango.log")
	megaMedia.Init("mega.log")

	filter.Init()

	manager.Init("manager.log")
	go manager.CommanderRoutine("124.232.133.211", 6379, "dcc-124.232.133.211")

	panic(http.ListenAndServe(":18124", nil))
}
