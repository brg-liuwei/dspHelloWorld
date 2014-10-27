package megaMedia

import (
	"io/ioutil"
	"net/http"

	proto "code.google.com/p/goprotobuf/proto"

	"common"
	"logger"
)

func NewBidRequest(r *http.Request) (req *MgxBidRequest) {
	if r.Method != "POST" {
		megaMediaLogger.Log(logger.ERROR, "megaMedia http.Request Method error: ", r.Method)
		return nil
	}

	req = new(MgxBidRequest)
	if body, err := ioutil.ReadAll(r.Body); err != nil {
		megaMediaLogger.Log(logger.ERROR, "megaMedia read http body error: ", err)
		return nil
	} else if err = proto.Unmarshal(body, req); err != nil {
		megaMediaLogger.Log(logger.ERROR, "megaMedia proto unmarshal error: ", err)
		return nil
	}

	if *req.IsTest == true {
		megaMediaLogger.Log(logger.DEBUG, "megaMedia test package")
	} else if *req.IsPing == true {
		megaMediaLogger.Log(logger.DEBUG, "megaMedia ping package")
	}
	return req
}

func (req *MgxBidRequest) ParseToCommon() (cb *common.BidRequest) {
	cb = new(common.BidRequest)
	cb.Id = req.GetBidGuid()
	if req.VisitorId != nil {
		// TODO: search cookie
		// cb.cookie = GetMegaCookieMapping(strconv.Itoa((int)(*cb.VisitorId)))
	}
	if req.DspUserId != nil {
		megaMediaLogger.Log(logger.INFO, "megaMedia why dsp user id != nil?: ", *req.DspUserId)
	}

	cb.Ip = req.GetIp()
	cb.Ua = req.GetUserAgent()
	// req.GetFlashVersion()
	cb.Language = req.GetLanguage()
	// req.UserSegment, 人群定向标签，暂不用
	// req.ContentTags, 内容标签，暂不用
	// req.SiteId, 在adx中的网站标识，暂不用
	cb.AdxId = common.MEGAMEDIA

	cb.Slots = make([]common.AdSlotType, len(req.Adslot))

	for i, s := range req.GetAdslot() {
		// megaMedia 需要回复广告位序列ID和广告位在adx中的唯一ID，拼接成ImpId存储
		cb.Slots[i].ImpId = ImpIdEncoding(s.GetId(), s.GetMegaxAid())
		//TODO: 根据买方集合ID来获取底价
		//cb.BidFloor =
		cb.Slots[i].W = int(s.GetAdsWidth())
		cb.Slots[i].H = int(s.GetAdsHeight())
		if s.GetPageNum() == 1 {
			cb.Slots[i].Visibility = 1
		} else {
			cb.Slots[i].Visibility = 0
		}
	}
	return
}
