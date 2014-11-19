package mango

import (
	"fmt"
	"net/http"

	"bid"
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
	bidRequest := NewBidRequest(r)
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
	commonResponse, _ := bid.Bid(commonRequest)
	fmt.Printf("\ncommon response: %#v\n", *commonResponse)
	bidResponse := new(BidResponse)
	bidResponse.ParseFromCommon(commonResponse)
	bidResponse.Response(w)
}
