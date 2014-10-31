package main

import (
	"adaptor/mango"
	"adaptor/megaMedia"
	"filter"
	"manager"

	//"fmt"
	"net/http"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	mango.Init("mango.log")
	megaMedia.Init("mega.log")
	filter.Init()

	manager.Init("manager.log")
	go manager.CommanderRoutine("124.232.133.211", 6379, "dcc-124.232.133.211")

	panic(http.ListenAndServe(":18124", nil))
}
