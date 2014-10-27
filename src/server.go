package main

import (
	"adaptor/mango"
	"adaptor/megaMedia"
	"bid"
	"common"
	"filter"
	"manager"

	"fmt"
	"net/http"
	"runtime"
)

func MangoClickHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n[CLICK] %#v", *r)
	// log to bid event
	w.Write([]byte("ok"))
}

func MangoDisplayHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n[WIN] %#v", *r)
	// log to bid event
	w.Write([]byte("ok"))
}

func MangoWinHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n[WIN] %#v", *r)
	// log to bid event
	w.Write([]byte("ok"))
}

func MangoBidHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\nbidrequest: \n", *r)
	bidRequest := mango.NewBidRequest(r)
	if bidRequest == nil {
		return
	}
	fmt.Printf("\nmango bidrequest: %#v\n", *bidRequest)

	//fmt.Println("mango bid request:")
	//fmt.Printf("%+v\n", *bidRequest)
	//fmt.Println()
	//fmt.Printf("%#v\n", *bidRequest)

	commonRequest := bidRequest.ParseToCommon()
	fmt.Printf("\ncommon bidrequest: %#v\n", *commonRequest)
	commonResponse := bid.Bid(commonRequest)
	fmt.Printf("\ncommon response: %#v\n", *commonResponse)
	bidResponse := new(mango.BidResponse)
	bidResponse.ParseFromCommon(commonResponse)
	bidResponse.Response(w)
}

func YeskyBidHandler(w http.ResponseWriter, r *http.Request) {
	bidRequest := megaMedia.NewBidRequest(r)
	if bidRequest == nil {
		return
	}
	fmt.Printf("\nyesky bidrequest: %#v\n", *bidRequest)

	commonRequest := bidRequest.ParseToCommon()
	fmt.Printf("\ncommon bidrequest: %#v\n", *commonRequest)

	commonResponse := bid.Bid(commonRequest)
	fmt.Printf("\ncommon response: %#v\n", *commonResponse)

	bidResponse := new(megaMedia.MgxBidResponse)
	bidResponse.ParseFromCommon(commonResponse)

	bidResponse.Response(w)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	manager.Init("manager.log")
	mango.Init("mango.log")
	filter.Init()

	// =========== debug ===========
	common.WinUrl = "124.232.133.211:18124/mango/win"
	// =========== debug ===========
	go manager.CommanderRoutine("124.232.133.211", 6379, "dcc-124.232.133.211")

	http.HandleFunc("/mango/bid", MangoBidHandler)
	http.HandleFunc("/mango/win", MangoWinHandler)
	http.HandleFunc("/mango/click", MangoClickHandler)
	http.HandleFunc("/mango/display", MangoDisplayHandler)

	http.HandleFunc("/yesky/bid", YeskyBidHandler)

	panic(http.ListenAndServe(":18124", nil))
}
