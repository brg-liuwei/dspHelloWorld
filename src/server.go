package main

import (
	"adaptor/mango"
	"logger"
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
	bidRequest := mango.NewBidRequest(r)
	if bidRequest == nil {
		return
	}

	//fmt.Println("mango bid request:")
	//fmt.Printf("%+v\n", *bidRequest)
	//fmt.Println()
	//fmt.Printf("%#v\n", *bidRequest)

	commonRequest := bidRequest.ParseToCommon()
	commonResponse := commonRequest.GenResponse()
	bidResponse := new(mango.BidResponse)
	bidResponse.ParseFromCommon(commonResponse)
	bidResponse.Response(w)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if ok := logger.Init("dspLog.log"); !ok {
		panic("logger init error.")
	}
	http.HandleFunc("/mango", MangoHandler)
	panic(http.ListenAndServe(":12306", nil))
}
