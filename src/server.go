package main

import (
	"adaptor/mango"
	"bid"
	"filter"
	"manager"

	"fmt"
	"net/http"
	"runtime"
)

func TestManager() {
	c := manager.NewCommand()
	jstr := `{
        "oper_type":"4",
        "fmt_ver":"1",
        "data":
        [   
        "adfaf2131dfafb",
        "adfaf2131dfafb",
        "adfaf2131dfafb"
        ]   
    }`
	if c.Parse(jstr) {
		fmt.Println("order_type: ", c.Ctype)
		fmt.Println("version: ", c.Cversion)
		fmt.Println("data: ", c.Data)
	} else {
		fmt.Println("Parse error")
	}
}

func MangoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("bidrequest: ", *r)
	bidRequest := mango.NewBidRequest(r)
	if bidRequest == nil {
		return
	}

	//fmt.Println("mango bid request:")
	//fmt.Printf("%+v\n", *bidRequest)
	//fmt.Println()
	//fmt.Printf("%#v\n", *bidRequest)

	commonRequest := bidRequest.ParseToCommon()
	commonResponse := bid.Bid(commonRequest)
	bidResponse := new(mango.BidResponse)
	bidResponse.ParseFromCommon(commonResponse)
	bidResponse.Response(w)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	manager.Init("manager.log")
	mango.Init("mango.log")
	filter.Init()

	//common.WinUrl = //...
	go manager.CommanderRoutine("<redis-ip>", 0 /* <redis port> */, "dcc-<dsp-ip>")
	http.HandleFunc("/mango", MangoHandler)
	panic(http.ListenAndServe(":18124", nil))
}
