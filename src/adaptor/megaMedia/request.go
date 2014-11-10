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
		megaMediaLogger.Log(logger.ERROR, "megaMedia http.Request Method error: ", r.Method, ", r: ", *r)
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

	if req.IsTest != nil && *req.IsTest == true {
		megaMediaLogger.Log(logger.DEBUG, "megaMedia test package")
	} else if req.IsPing != nil && *req.IsPing == true {
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
		megaMediaLogger.Log(logger.INFO, "megaMedia why dsp user id != nil?: ", req.GetDspUserId())
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

		if len(s.BuyerId) != len(s.BuyerMinPrice) {
			megaMediaLogger.Log(logger.ERROR, "len(Id) isnot equal to len(MinPrice)")
		} else {
			for j, buyer := range s.BuyerId {
				if buyer == "331080" {
					cb.Slots[i].BidFloor = int(s.BuyerMinPrice[j] * 1000)
					break
				}
			}
		}

		cb.Slots[i].W = int(s.GetAdsWidth())
		cb.Slots[i].H = int(s.GetAdsHeight())
		if s.GetPageNum() == 1 {
			cb.Slots[i].Visibility = 1
		} else {
			cb.Slots[i].Visibility = 0
		}

		creatives := s.GetCreativeFiles()
		if len(creatives) > 0 {
			cb.Slots[i].CreativeType = make([]common.AdType, len(creatives))
			for j, creative := range creatives {
				switch {
				// jpg, png, gif
				case creative == "1" || creative == "2" || creative == "3":
					cb.Slots[i].CreativeType[j] = common.Banner
				case creative == "4": // swf
					cb.Slots[i].CreativeType[j] = common.Flash
					// flv, mp4
				case creative == "5" || creative == "6":
					cb.Slots[i].CreativeType[j] = common.Video
					// html, html5
				case creative == "7" || creative == "8":
					cb.Slots[i].CreativeType[j] = common.Html
				default:
					cb.Slots[i].CreativeType[j] = common.AdtypeUnknown
				}
			}
		}

	}
	return
}
