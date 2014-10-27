package mango

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"common"
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
		mangoLogger.Log(logger.ERROR, "mango http.Request Method error: ", r.Method)
		return nil
	}

	var m map[string]interface{}
	if body, err := ioutil.ReadAll(r.Body); err != nil {
		mangoLogger.Log(logger.ERROR, "mango read http.Request.Body error: ", err)
		return nil
	} else {
		m = make(map[string]interface{})
		if err := json.Unmarshal(body, &m); err != nil {
			mangoLogger.Log(logger.ERROR, "unmarshal json error: ", err)
			return nil
		}
	}

	br := new(BidRequest)
	if !br.SetId(&m) {
		mangoLogger.Log(logger.ERROR, "bid request set id error")
		return nil
	}
	if !br.SetImps(&m) {
		mangoLogger.Log(logger.ERROR, "bid request set imps error")
		return nil
	}
	if !br.SetApp(&m) {
		mangoLogger.Log(logger.ERROR, "bid request set app error")
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

func (r *BidRequest) ParseToCommon() (cb *common.BidRequest) {
	cb = new(common.BidRequest)
	cb.Id = r.Id
	cb.Mobile = true
	cb.Ip = r.Dev.Ip
	cb.Did = r.Dev.Did
	cb.Dpid = r.Dev.Dpid
	cb.Mac = r.Dev.Mac
	/* TODO:
	   cb.Region =
	   cb.Uuid =
	*/
	cb.Ua = r.Dev.Ua
	cb.Os = r.Dev.Os
	/* TODO: (手机端，可能没有Browser,Url)
	   cb.Browser =
	   cb.Url =
	   cb.Domain =
	*/
	cb.Language = r.Dev.Language
	cb.Media = "mango"

	/* mango的Imps数组长度固定为1 */
	cb.Slots = make([]common.AdSlotType, 0, 1)
	var slot common.AdSlotType
	if len(r.Imps) == 0 {
		return
	}
	imp := r.Imps[0]
	slot.ImpId = imp.ImpId
	slot.BidFloor = imp.BidFloor
	slot.W = imp.W
	slot.H = imp.H
	switch imp.Pos {
	case TOPVIEW:
		slot.Visibility = 20
	case BOTTOMVIEW:
		slot.Visibility = 15
	case TOPROLL:
		slot.Visibility = 10
	case BOTTOMROLL:
		slot.Visibility = 5
	default:
		slot.Visibility = 1
	}
	slot.CatOut = make(map[string]bool)
	for _, adType := range imp.Btype {
		slot.CatOut[string(adType)] = true
	}
	slot.AttrOut = make(map[string]bool)
	for _, adAttr := range imp.Battr {
		slot.AttrOut[string(adAttr)] = true
	}
	slot.Instl = imp.Instl
	slot.Splash = imp.Splash
	cb.Slots = append(cb.Slots, slot)

	cb.AdxId = common.MANGO
	// cb.User = ... no user info
	// cb.Site = ... no site info
	cb.App.Id = r.App.Aid
	// TODO: calc cb.App.Quality
	return
}
