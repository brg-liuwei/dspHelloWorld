package mango

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"logger"
)

type BidRequest struct {
	Id   string       // bid id
	Imps []Impression // Array of Impression
	App  Application  // Application Info
	Dev  Device       // Device Info
	Bcat []AdCategory // black list of ad category
	Badv []string     // blask list of ad url
}

func NewBidRequest(r *http.Request) *BidRequest {
	if r.Method != "POST" {
		logger.LOG(logger.ERROR, "http.Request Method error: ", r.Method)
		return nil
	}

	var m map[string]interface{}
	if body, err := ioutil.ReadAll(r.Body); err != nil {
		logger.LOG(logger.ERROR, "read ttp.Request.Body error: ", err)
		return nil
	} else {
		m = make(map[string]interface{})
		if err := json.Unmarshal(body, &m); err != nil {
			logger.LOG(logger.ERROR, "unmarshal json error: ", err)
			return nil
		}
	}

	br := new(BidRequest)
	if !br.SetId(&m) {
		logger.LOG(logger.ERROR, "bid request set id error")
		return nil
	}
	if !br.SetImps(&m) {
		logger.LOG(logger.ERROR, "bid request set imps error")
		return nil
	}
	if !br.SetApp(&m) {
		logger.LOG(logger.ERROR, "bid request set app error")
		return nil
	}
	br.SetDevice()
	br.SetBcat()
	br.SetBadv()
	return br
}
