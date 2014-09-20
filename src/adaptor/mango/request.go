package mango

import (
	"encoding/json"
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
		logger.Log(logger.ERROR, "http.Request Method error: ", r.Method)
		return nil
	}

	var m map[string]interface{}
	if body, err := ioutil.ReadAll(r.Body); err != nil {
		logger.Log(logger.ERROR, "read ttp.Request.Body error: ", err)
		return nil
	} else {
		m = make(map[string]interface{})
		if err := json.Unmarshal(body, &m); err != nil {
			logger.Log(logger.ERROR, "unmarshal json error: ", err)
			return nil
		}
	}

	br := new(BidRequest)
	if !br.SetId(&m) {
		logger.Log(logger.ERROR, "bid request set id error")
		return nil
	}
	if !br.SetImps(&m) {
		logger.Log(logger.ERROR, "bid request set imps error")
		return nil
	}
	if !br.SetApp(&m) {
		logger.Log(logger.ERROR, "bid request set app error")
		return nil
	}
	br.SetDevice(&m)
	br.SetBcat(&m)
	br.SetBadv(&m)
	return br
}

func (r *BidRequest) SetId(m *map[string]interface{}) bool {
	if v, ok := (*m)["id"]; !ok {
		return false
	} else if r.Id, ok = v.(string); !ok {
		return false
	} else {
		return true
	}
}

func (r *BidRequest) SetImps(m *map[string]interface{}) bool {
	var imps []interface{}
	if v, ok := (*m)["imp"]; !ok {
		return false
	} else if imps, ok = v.([]interface{}); !ok {
		return false
	} else {
		r.Imps = make([]Impression, len(imps))
		for i := 0; i != len(imps); i++ {
			if imp, ok := imps[i].(map[string]interface{}); !ok {
				return false
			} else if !r.Imps[i].Assign(&imp) {
				return false
			}
		}
		return true
	}
}

func (r *BidRequest) SetApp(m *map[string]interface{}) bool {
	var appMap map[string]interface{}
	if v, ok := (*m)["app"]; !ok {
		return false
	} else if appMap, ok = v.(map[string]interface{}); !ok {
		return false
	} else if !r.App.Assign(&appMap) {
		return false
	} else {
		return true
	}
}

func (r *BidRequest) SetDevice(m *map[string]interface{}) {
	var devMap map[string]interface{}
	if v, ok := (*m)["device"]; ok {
		if devMap, ok = v.(map[string]interface{}); ok {
			r.Dev.Assign(&devMap)
		}
	}
}

func (r *BidRequest) SetBcat(m *map[string]interface{}) {
	if v, ok := (*m)["bcat"]; ok {
		if array, ok := v.([]interface{}); ok {
			r.Bcat = make([]AdCategory, 0, len(array))
			for _, elem := range array {
				if e, ok := elem.(string); ok {
					r.Bcat = append(r.Bcat, AdCategory(e))
				}
			}
		}
	}
}

func (r *BidRequest) SetBadv(m *map[string]interface{}) {
	if v, ok := (*m)["badv"]; ok {
		if array, ok := v.([]interface{}); ok {
			r.Badv = make([]string, 0, len(array))
			for _, elem := range array {
				if e, ok := elem.(string); ok {
					r.Badv = append(r.Badv, e)
				}
			}
		}
	}
}
