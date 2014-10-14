package mango

import (
	"encoding/json"
	"fmt"
	"net/http"

	"common"
	"logger"
)

type BidResponse struct {
	Id      string // 对应相应的BidRequest id
	Bidid   string // dsp自己的竞价ID
	Nbr     int    // 不参与竞价的原因(默认为0，表示未知),参与竞价时json不包含此域
	Seatbid []SeatBid
}

type SeatBid struct {
	Bid  []Bid
	Seat string // Seat Bid Id, 由dsp生成
}

type Bid struct {
	Impid   string
	Price   int // unit: 0.0001RMB/CPM, eg: 0.6RMB=6000
	Adid    string
	Nurl    string // win notice url
	Adm     string // HTML片段或者url地址
	Adw     int
	Adh     int
	Iurl    string   // display monitor url(当Adm中不包含时使用)
	Curl    string   // Adm的#CLICK_URL#宏，被点击的跳转地址
	Cturl   []string // click monitor url(用于Adm中的#CLICK_URL#宏),URL需要支持重定向
	Cid     string   // Campaign ID
	Crid    string   // 物料ID
	Ctype   ClickType
	Cbundle string // 当Ctype为DownloadApp时，为下载app的包名
	Attr    []AdAttr
	Domain  string // 广告主域名
	Ext     External
}

type External struct {
	Instl int    // 全插屏广告需返回: 1-插屏，2-全屏
	Adt   int    // 开屏广告必须返回, 播放时长, 单位: 秒
	Ade   string // 开屏广告必须返回, 过期时间, 格式: YYYYMMdd
}

func (bid *Bid) fillMap(m *map[string]interface{}) {
	(*m)["impid"] = bid.Impid
	(*m)["price"] = bid.Price
	(*m)["adid"] = bid.Adid
	if len(bid.Nurl) != 0 {
		(*m)["nurl"] = bid.Nurl
	}
	(*m)["adm"] = bid.Adm
	if bid.Adw != 0 && bid.Adh != 0 {
		(*m)["adw"] = bid.Adw
		(*m)["adh"] = bid.Adh
	}
	if len(bid.Iurl) != 0 {
		(*m)["iurl"] = bid.Iurl
	}
	(*m)["curl"] = bid.Curl
	if len(bid.Cturl) != 0 {
		(*m)["cturl"] = bid.Cturl
	}
	if len(bid.Cid) != 0 {
		(*m)["cid"] = bid.Cid
	}
	if len(bid.Crid) != 0 {
		(*m)["crid"] = bid.Crid
	}
	if bid.Ctype != CTypeUnkown {
		(*m)["ctype"] = bid.Ctype
		if bid.Ctype == DownloadApp {
			(*m)["cbundle"] = bid.Cbundle
		}
	}
	if len(bid.Attr) != 0 {
		(*m)["attr"] = bid.Attr
	}
	if len(bid.Domain) != 0 {
		(*m)["adomain"] = bid.Domain
	}

	/* ext */
	extMap := make(map[string]interface{})
	if bid.Ext.Instl != 0 {
		extMap["instl"] = bid.Ext.Instl
	}
	if bid.Ext.Adt != 0 {
		extMap["adt"] = bid.Ext.Adt
	}
	if len(bid.Ext.Ade) != 0 {
		extMap["ade"] = bid.Ext.Ade
	}
	(*m)["ext"] = extMap
}

func (r *BidResponse) Response(w http.ResponseWriter) {
	var sbSlice []map[string]interface{}
	rmap := make(map[string]interface{})
	rmap["id"] = r.Id
	if len(r.Bidid) != 0 {
		rmap["bidid"] = r.Bidid
	}
	if len(r.Seatbid) == 0 { /* give up */
		rmap["nbr"] = r.Nbr
		goto end
	}

	/* for seat bid */
	sbSlice = make([]map[string]interface{}, 0, len(r.Seatbid))
	for i := 0; i != len(r.Seatbid); i++ {
		sb := make(map[string]interface{})

		// seat
		if len(r.Seatbid[i].Seat) != 0 {
			sb["seat"] = r.Seatbid[i].Seat
		}

		// bid
		bidSlice := make([]map[string]interface{}, 0, len(r.Seatbid[i].Bid))
		for j := 0; j != len(r.Seatbid[i].Bid); j++ {
			bidMap := make(map[string]interface{})
			bid := r.Seatbid[i].Bid[j]
			bid.fillMap(&bidMap)
			bidSlice = append(bidSlice, bidMap)
		}
		sb["bid"] = bidSlice
		sbSlice = append(sbSlice, sb)
	}
	rmap["seatbid"] = sbSlice

end:
	response, _ := json.Marshal(rmap)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(response)))
	w.Write(response)

	/* log bid response */
	mangoLogger.Log(logger.INFO, "{\"MONGO\":"+string(response)+"}")
}

func (r *BidResponse) ParseFromCommon(cr *common.BidResponse) {
	r.Id = cr.ReqId
	r.Bidid = cr.BidId
	if len(cr.Ads) == 0 {
		r.Nbr = 0
		return
	}
	// mongo只返回1个seatbid
	r.Seatbid = make([]SeatBid, 0, 1)
	var sb SeatBid
	sb.Bid = make([]Bid, 0, 1)

	commonBid := cr.Ads[0]
	var mongoBid Bid
	mongoBid.Impid = commonBid.ImpId
	mongoBid.Price = commonBid.Price
	mongoBid.Adid = commonBid.AdId
	mongoBid.Nurl = commonBid.WinUrl
	mongoBid.Adm = commonBid.Adm
	mongoBid.Adw = commonBid.W
	mongoBid.Adh = commonBid.H
	mongoBid.Iurl = commonBid.DisplayMonitor
	mongoBid.Curl = commonBid.LandingPage
	mongoBid.Cturl = make([]string, 0)
	if len(commonBid.ClickMonitor) != 0 {
		mongoBid.Cturl = append(mongoBid.Cturl, commonBid.ClickMonitor)
	}
	/* TODO:
	   mongoBid.Cid
	   mongoBid.Crid
	   mongoBid.Cbundle
	   mongoBid.Attr
	*/
	mongoBid.Domain = commonBid.Domain

	// TODO: Ext
	mongoBid.Ext.Instl = 1        // 1-插屏，2-全屏
	mongoBid.Ext.Adt = 5          // 开屏播放时长: 5s
	mongoBid.Ext.Ade = "20150101" // 开屏过期时间: YYYYMMdd

	sb.Bid = append(sb.Bid, mongoBid)
	r.Seatbid = append(r.Seatbid, sb)
}
